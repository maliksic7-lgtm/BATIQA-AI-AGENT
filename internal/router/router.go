package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"batiqa-ai/internal/events"
	"batiqa-ai/internal/handler"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
	ticketservice "batiqa-ai/internal/service/ticket"
)

// New creates the HTTP router for Phase 1 (health only, no DB).
// Kept for backward compatibility and tests without DB.
func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handler.HealthCheck)
	mux.Handle("/api/docs/", handler.NewOpenAPIHandler())

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" || strings.HasPrefix(r.URL.Path, "/api/docs/") {
			mux.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			handler.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
			return
		}
		// Try static for non-API
		if serveStatic(w, r) {
			return
		}
		handler.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found")
	})

	return chain(wrapped, recoveryMiddleware, loggingMiddleware, corsMiddleware)
}

// NewWithDB creates router with Phase 4 ticket & chat endpoints plus Phase 5 hotel-info.
// Flow: Handler -> Service -> Repository -> MongoDB per DATABASE.md
func NewWithDB(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	// Health (DB-aware when a database is present)
	if db != nil {
		mux.HandleFunc("/api/health", handler.NewHealthHandlerWithDB(db).ServeHTTP)
	} else {
		mux.HandleFunc("/api/health", handler.HealthCheck)
	}

	// OpenAPI / Swagger UI (works even without DB)
	mux.Handle("/api/docs/", handler.NewOpenAPIHandler())

	if db != nil {
		// Repositories
		ticketRepo := repository.NewTicketRepository(db)
		guestRepo := repository.NewGuestRepository(db)
		convRepo := repository.NewConversationRepository(db)
		hotelRepo := repository.NewHotelInfoRepository(db)
		recRepo := repository.NewRecommendationRepository(db)
		sessionRepo := repository.NewStaffSessionRepository(db)
		orderRepo := repository.NewRestaurantOrderRepository(db)

		// Services
		broker := events.New()
		aiSvc := ai.NewService()
		ticketSvc := ticketservice.NewService(ticketRepo, guestRepo)
		ticketSvc.SetBroker(broker)

		// Guests may only open their own tickets: DB-backed ownership check.
		handler.SetTicketOwnershipChecker(func(r *http.Request, ticketNumber, room string) bool {
			t, err := ticketSvc.GetByTicketNumber(ticketNumber)
			return err == nil && t != nil && strings.EqualFold(t.RoomNumber, room)
		})

		// Handlers
		ticketHandler := handler.NewTicketHandler(ticketSvc)
		chatHandler := handler.NewChatHandlerFull(aiSvc, ticketSvc, convRepo, guestRepo, hotelRepo, recRepo)
		hotelHandler := handler.NewHotelHandler(hotelRepo)
		recHandler := handler.NewRecommendationHandler(recRepo)
		staffAuthHandler := handler.NewStaffAuthHandler(repository.NewStaffRepository(db), sessionRepo)
		// Production login throttling: persistent across restarts & instances.
		handler.SetRateLimiter(handler.NewMongoRateLimiter(repository.NewLoginRateLimitRepository(db)))
		statsHandler := handler.NewStatsHandler(ticketRepo)
		staffRepo := repository.NewStaffRepository(db)
		assignHandler := handler.NewAssignHandler(ticketSvc, staffRepo, repository.NewAssignmentRepository(db))
		convHandler := handler.NewConversationHandler(convRepo)
		qrHandler := handler.NewQRHandler(staffAuthHandler)
		sseHandler := handler.NewSSEHandler(staffAuthHandler, broker)
		analyticsHandler := handler.NewAnalyticsHandler(ticketRepo)
		infographicsHandler := handler.NewInfographicsHandlerFull(ticketRepo, convRepo, orderRepo)
		restaurantHandler := handler.NewRestaurantHandler(orderRepo)

		// Staff auth: POST /api/staff/login (public), GET /api/staff/me, POST /api/staff/logout (auth)
		mux.HandleFunc("/api/staff/login", staffAuthHandler.Login)
		mux.Handle("/api/staff/me", staffAuthHandler.AuthMiddleware(http.HandlerFunc(staffAuthHandler.Me)))
		mux.Handle("/api/staff/logout", staffAuthHandler.AuthMiddleware(http.HandlerFunc(staffAuthHandler.Logout)))

		// Stats: GET /api/tickets/stats (staff only)
		mux.Handle("/api/tickets/stats", staffAuthHandler.AuthMiddleware(http.HandlerFunc(statsHandler.ServeHTTP)))

		// Live updates: GET /api/events (staff, SSE)
		mux.Handle("/api/events", http.HandlerFunc(sseHandler.ServeHTTP))

		// QR room codes: GET /api/rooms/{room}/qr (staff only)
		mux.Handle("/api/rooms/", staffAuthHandler.AuthMiddleware(http.HandlerFunc(qrHandler.ServeHTTP)))

		// Analytics: GET /api/analytics (staff only)
		mux.Handle("/api/analytics", staffAuthHandler.AuthMiddleware(http.HandlerFunc(analyticsHandler.ServeHTTP)))

		// Infographics: GET /api/analytics/infographics (staff only)
		mux.Handle("/api/analytics/infographics", staffAuthHandler.AuthMiddleware(http.HandlerFunc(infographicsHandler.ServeHTTP)))

		// Restaurant: public menu; guest places+lists own orders; staff lists all +
		// updates status.
		mux.HandleFunc("/api/menu", handler.MenuHandler)
		mux.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/status") && r.Method == http.MethodPatch {
				staffAuthHandler.AuthMiddleware(http.HandlerFunc(restaurantHandler.UpdateStatus)).ServeHTTP(w, r)
				return
			}
			handler.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
		})
		mux.Handle("/api/orders", staffAuthHandler.EitherAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				restaurantHandler.Create(w, r)
			case http.MethodGet:
				restaurantHandler.List(w, r)
			default:
				handler.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			}
		})))

		// Guest identity from QR token: GET /api/guest/me
		mux.Handle("/api/guest/me", handler.GuestAuthMiddleware(http.HandlerFunc(chatHandler.GuestMe)))

		// Chat text + photo: guests must present a valid QR token
		mux.Handle("/api/chat/photo", handler.GuestAuthMiddleware(http.HandlerFunc(chatHandler.Photo)))
		mux.Handle("/api/chat", handler.GuestAuthMiddleware(http.HandlerFunc(chatHandler.ServeHTTP)))

		// Chat history per session (guest or staff)
		mux.Handle("/api/conversations", staffAuthHandler.EitherAuth(http.HandlerFunc(convHandler.ServeHTTP)))

		// Hotel info & recommendations stay public (non-guest-specific data)
		mux.HandleFunc("/api/hotel-info", hotelHandler.ServeHTTP)
		mux.HandleFunc("/api/hotel_info", hotelHandler.ServeHTTP)
		mux.HandleFunc("/api/recommendations", recHandler.ServeHTTP)

		// Tickets: POST /api/tickets, GET /api/tickets — staff OR scoped guest
		mux.Handle("/api/tickets", staffAuthHandler.EitherAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				ticketHandler.Create(w, r)
			case http.MethodGet:
				ticketHandler.List(w, r)
			default:
				handler.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			}
		})))

		// Tickets detail and sub-resources: /api/tickets/{id}, /{id}/status (auth),
		// /{id}/priority (auth), /{id}/assign + /{id}/assignments (auth for POST)
		mux.HandleFunc("/api/tickets/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			switch {
			case strings.HasSuffix(path, "/status"):
				staffAuthHandler.AuthMiddleware(http.HandlerFunc(ticketHandler.UpdateStatus)).ServeHTTP(w, r)
			case strings.HasSuffix(path, "/priority"):
				staffAuthHandler.AuthMiddleware(http.HandlerFunc(ticketHandler.UpdatePriority)).ServeHTTP(w, r)
			case strings.HasSuffix(path, "/assign") && r.Method == http.MethodPost:
				staffAuthHandler.AuthMiddleware(http.HandlerFunc(assignHandler.ServeHTTP)).ServeHTTP(w, r)
			case strings.HasSuffix(path, "/assignments") || strings.HasSuffix(path, "/assign"):
				staffAuthHandler.EitherAuth(http.HandlerFunc(assignHandler.ServeHTTP)).ServeHTTP(w, r)
			case r.Method == http.MethodGet:
				// Guests may only read their own room's tickets; staff any.
				if handler.GuestRoom(r) != "" {
					handler.GuestScopedDetail(http.HandlerFunc(ticketHandler.GetDetail)).ServeHTTP(w, r)
					return
				}
				ticketHandler.GetDetail(w, r)
			default:
				handler.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			}
		})
	}

	// Wrapped handler: API routes via mux, static files for frontend, 503 if DB nil
	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API routes via mux
		apiPrefixes := []string{
			"/api/health", "/api/chat", "/api/conversations", "/api/tickets",
			"/api/hotel-info", "/api/hotel_info", "/api/recommendations",
			"/api/staff/login", "/api/staff/me", "/api/staff/logout",
			"/api/guest/me", "/api/rooms/", "/api/events", "/api/analytics",
			"/api/menu", "/api/orders",
			"/api/docs",
		}
		for _, p := range apiPrefixes {
			base := strings.TrimSuffix(p, "/")
			if r.URL.Path == base || strings.HasPrefix(r.URL.Path, base+"/") {
				// DB required for all except health
				if db == nil && p != "/api/health" && p != "/api/docs" {
					handler.WriteError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database not available")
					return
				}
				mux.ServeHTTP(w, r)
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			handler.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
			return
		}
		// Static files for guest interface (Phase 5)
		if serveStatic(w, r) {
			return
		}
		handler.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found")
	})

	return chain(wrapped, recoveryMiddleware, loggingMiddleware, corsMiddleware)
}

// serveStatic tries to serve web assets for Phase 5 guest interface.
// Returns true if served (or 404 handled), false if not a static path.
func serveStatic(w http.ResponseWriter, r *http.Request) bool {
	webDir := findWebDir()
	if webDir == "" {
		return false
	}

	// Clean path
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	// Map URL path to filesystem
	// / -> web/index.html, /guest/* -> web/guest/*, /css/* -> web/css/*, /js/* -> web/js/*, /guest -> web/guest/index.html
	fullPath := filepath.Join(webDir, filepath.FromSlash(path))

	// Handle directory -> try index.html
	if strings.HasSuffix(path, "/") {
		fullPath = filepath.Join(fullPath, "index.html")
	}
	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		// Try with .html extension for clean URLs? e.g., /guest -> /guest/index.html
		if path == "/guest" {
			alt := filepath.Join(webDir, "guest", "index.html")
			if _, err2 := os.Stat(alt); err2 == nil {
				http.ServeFile(w, r, alt)
				return true
			}
		}
		return false
	}
	if info.IsDir() {
		// Serve index.html inside
		index := filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(index); err == nil {
			http.ServeFile(w, r, index)
			return true
		}
		return false
	}
	http.ServeFile(w, r, fullPath)
	return true
}

func findWebDir() string {
	candidates := []string{
		"web",
		"../web",
		"../../web",
		filepath.Join(filepath.Dir(os.Args[0]), "web"),
		filepath.Join(filepath.Dir(os.Args[0]), "../web"),
	}
	// Also try exe dir
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "web"))
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "..", "web"))
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "..", "..", "web"))
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			// Check has index.html
			if _, err := os.Stat(filepath.Join(c, "index.html")); err == nil {
				return c
			}
		}
	}
	return ""
}

// chain applies middlewares in order (first is outermost)
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
