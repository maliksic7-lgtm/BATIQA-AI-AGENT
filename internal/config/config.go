package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds all environment-based configuration.
// Phase 2 adds MySQL database settings.
type Config struct {
	Port               string
	Env                string
	LogLevel           string
	DBDSN              string
	GeminiKey          string
	DBMaxOpen          int
	DBMaxIdle          int
	CORSAllowedOrigins []string
}

// Load reads configuration from environment variables with sensible defaults.
// It also attempts to load a .env file by walking up from the working directory (simple parser, no external dep).
func Load() *Config {
	// Try to load .env file silently (best-effort, no error if missing)
	loadDotEnvUpward(".env")

	cfg := &Config{
		Port:               envOrDefault("PORT", "8080"),
		Env:                envOrDefault("ENV", "development"),
		LogLevel:           envOrDefault("LOG_LEVEL", "info"),
		DBDSN:              os.Getenv("DB_DSN"),
		GeminiKey:          os.Getenv("GEMINI_API_KEY"),
		DBMaxOpen:          envIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdle:          envIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		CORSAllowedOrigins: parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}
	// Default DBDSN for local dev only if not set (empty password, not hardcoded secret)
	if cfg.DBDSN == "" {
		cfg.DBDSN = "root:@tcp(127.0.0.1:3306)/batiqa_ai?parseTime=true&charset=utf8mb4&loc=Local"
	}
	return cfg
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envIntOrDefault(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// IsProduction returns true when ENV=production
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.Env, "production")
}

// Addr returns HTTP listen address (e.g., ":8080")
func (c *Config) Addr() string {
	if strings.HasPrefix(c.Port, ":") {
		return c.Port
	}
	return ":" + c.Port
}

// parseOrigins splits a comma-separated origin list. Empty -> localhost defaults.
func parseOrigins(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return []string{"http://localhost:8080", "http://localhost:3000", "http://127.0.0.1:8080"}
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// loadDotEnvUpward walks up from the working directory (max 4 levels) to find the .env file.
func loadDotEnvUpward(name string) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			_ = loadDotEnv(path)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

// loadDotEnv is a minimal .env parser without external dependencies.
// It reads KEY=VALUE lines, ignores comments (#) and empty lines.
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Remove inline comments not inside quotes (simple)
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove surrounding quotes
		val = strings.Trim(val, `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}
