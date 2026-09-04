// BATIQA AI Guest Assistant - Service Worker
// Provides offline caching (app shell), runtime caching for public API,
// and notification handling for ticket status updates.
const VERSION = 'batiqa-v1.0.0';
const APP_SHELL_CACHE = VERSION + '-app';
const API_CACHE = VERSION + '-api';

// Static assets to precache for offline use (app shell).
const APP_SHELL = [
  '/',
  '/index.html',
  '/manifest.json',
  '/icon.svg',
  '/css/guest.css',
  '/css/staff.css',
  '/js/api.js',
  '/js/chat.js',
  '/js/staff.js',
  '/js/requests.js',
  '/js/info.js',
  '/guest/index.html',
  '/guest/chat.html',
  '/guest/info.html',
  '/guest/requests.html',
  '/staff/login.html',
  '/staff/index.html'
];

// Cache-first list. We only cache public, non-authenticated, 
// safely-reusable GET endpoints.
const API_CACHE_PATHS = [
  '/api/hotel-info',
  '/api/hotel_info',
  '/api/recommendations',
  '/api/health'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(APP_SHELL_CACHE)
      .then((cache) => cache.addAll(APP_SHELL))
      .then(() => self.skipWaiting())
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((keys) => Promise.all(
        keys.filter((key) => key !== APP_SHELL_CACHE && key !== API_CACHE)
            .map((key) => caches.delete(key))
      ))
      .then(() => self.clients.claim())
  );
});

// Helper: is this a navigation request?
function isNavigation(req) {
  return req.mode === 'navigate';
}

// Helper: should we cache this API URL?
function isCacheableApi(url) {
  return API_CACHE_PATHS.some((p) => url.pathname === p || url.pathname.startsWith(p + '/'));
}

// Fetch handler: cache-first for app shell & public API, network-first for
// everything else (with offline fallback to cached shell for navigations).
self.addEventListener('fetch', (event) => {
  const req = event.request;
  const url = new URL(req.url);

  // Only handle same-origin GET requests.
  if (url.origin !== self.location.origin || req.method !== 'GET') return;

  // Public API: stale-while-revalidate (fast + update in background).
  if (url.pathname.startsWith('/api/')) {
    if (isCacheableApi(url)) {
      event.respondWith(staleWhileRevalidate(req, url));
    }
    return;
  }

  // App shell (static assets): cache-first.
  // Exclude admin/content-prohibited paths from offline caching.
  if (isNavigation(req)) {
    event.respondWith(networkFirstForNav(req, url));
    return;
  }

  event.respondWith(
    caches.match(req).then((cached) => {
      if (cached) return cached;
      return fetch(req).then((res) => {
        if (res && res.ok && isCacheableStatic(url)) {
          const copy = res.clone();
          caches.open(APP_SHELL_CACHE).then((cache) => cache.put(req, copy));
        }
        return res;
      });
    })
  );
});

function isCacheableStatic(url) {
  return /\.(css|js|svg|png|jpg|jpeg|webp|ico|woff2?|json|html)$/i.test(url.pathname);
}

// Network-first for navigations so fresh HTML always loads; fall back to
// cached shell when offline.
function networkFirstForNav(req, url) {
  return fetch(req).then((res) => {
    if (res && res.ok) {
      const copy = res.clone();
      caches.open(APP_SHELL_CACHE).then((cache) => cache.put('/index.html', copy));
    }
    return res;
  }).catch(() => {
    // Try cached app shell (best-effort single page).
    if (url.pathname !== '/') {
      return caches.open(APP_SHELL_CACHE).then((c) => c.match(url));
    }
    return caches.open(APP_SHELL_CACHE).then((c) => c.match('/index.html'));
  });
}

// Stale-while-revalidate: serve cached copy immediately, refresh cache async.
function staleWhileRevalidate(req, url) {
  const cacheKey = req;
  return caches.open(API_CACHE).then((cache) =>
    cache.match(cacheKey).then((cached) => {
      const networkFetch = fetch(req).then((res) => {
        if (res && res.ok) cache.put(cacheKey, res.clone());
        return res;
      }).catch(() => cached);
      return cached || networkFetch;
    })
  );
}

// ---------- Push notifications ----------
// Handle background push messages from a push server (future). For now we
// also accept a "ticket-update" payload to surface a system notification.
self.addEventListener('push', (event) => {
  let data = {};
  try { data = event.data ? event.data.json() : {}; } catch (e) { /* ignore */ }

  const title = data.title || 'BATIQA AI';
  const options = {
    body: data.body || 'Pembaruan status permintaan Anda.',
    icon: '/icon.svg',
    badge: '/icon.svg',
    data: { url: data.url || '/guest/requests.html' }
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  const url = (event.notification.data && event.notification.data.url) || '/';
  event.waitUntil(
    self.clients.matchAll({ type: 'window', includeUncontrolled: true })
      .then((clientList) => {
        for (const client of clientList) {
          if ('focus' in client) { client.navigate(url); return client.focus(); }
        }
        return self.clients.openWindow(url);
      })
  );
});

self.addEventListener('message', (event) => {
  // 'SKIP_WAITING' lets an updated SW take control immediately.
  if (event.data === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});
