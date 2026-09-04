// BATIQA PWA bootstrap: service worker registration + notification helpers.
(function () {
  const SW_PATH = '/sw.js';

  function registerSW() {
    if (!('serviceWorker' in navigator)) return;

    window.addEventListener('load', function () {
      navigator.serviceWorker.register(SW_PATH)
        .then(function (reg) {
          // Update available -> tell SW to skip waiting so new version applies.
          reg.addEventListener('updatefound', function () {
            const newWorker = reg.installing;
            if (!newWorker) return;
            newWorker.addEventListener('statechange', function () {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                newWorker.postMessage('SKIP_WAITING');
              }
            });
          });
        })
        .catch(function () { /* SW optional for non-PWA-capable browsers */ });
    });
  }

  // Request notification permission and return a Notification if granted.
  // Only meaningful on secure contexts (https or localhost).
  function notify(title, opts) {
    try {
      if (!('Notification' in window)) return Promise.resolve(null);
      if (Notification.permission === 'denied') return Promise.resolve(null);
      if (Notification.permission === 'granted') {
        const n = new Notification(title, Object.assign({ icon: '/icon.svg', badge: '/icon.svg' }, opts));
        n.onclick = function () {
          if (opts && opts.url) window.location.href = opts.url;
          n.close();
        };
        return Promise.resolve(n);
      }
      // Not requested yet: ask once, then fire if granted.
      return Notification.requestPermission().then(function (p) {
        if (p === 'granted' && typeof Notification !== 'undefined') {
          return new Notification(title, Object.assign({ icon: '/icon.svg', badge: '/icon.svg' }, opts));
        }
        return null;
      });
    } catch (e) {
      return Promise.resolve(null);
    }
  }

  // Poll-driven foreground notification: fire a notification when any of the
  // given ticket statuses changes. Used on guest requests page as a lightweight
  // real-time cue (falls back gracefully if notifications unsupported).
  function watchTicketStatus(pollMs, onStatusES) {
    const INTERVAL = pollMs || 30000;
    const statusCache = {};
    const refresh = function (tickets) {
      (tickets || []).forEach(function (t) {
        const key = t.ticket_number || t._id;
        if (!key) return;
        const prev = statusCache[key];
        statusCache[key] = t.status;
        if (prev && prev !== t.status) {
          notify('Pembaruan Tiket', {
            body: 'Status ' + key + ' berubah: ' + prev + ' → ' + t.status,
            url: '/guest/requests.html'
          });
        }
      });
    };
    return {
      refresh: refresh,
      interval: INTERVAL
    };
  }

  // Expose helpers on window.
  window.BATIQA = window.BATIQA || {};
  window.BATIQA.notify = notify;
  window.BATIQA.watchTicketStatus = watchTicketStatus;

  registerSW();
})();
