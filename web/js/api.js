// BATIQA API Client - lightweight, no heavy libs
const API = {
  base: '',
  async request(path, opts = {}) {
    const res = await fetch(path, {
      headers: { 'Content-Type': 'application/json', ...guestHeaders(), ...(opts.headers || {}) },
      ...opts,
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      const msg = data?.error?.message || `Request failed ${res.status}`;
      const err = new Error(msg);
      err.status = res.status;
      err.data = data;
      throw err;
    }
    return data;
  },
  chat(session_id, room_number, message) {
    const body = { session_id, room_number, message };
    const t = getToken();
    if (t) body.token = t;
    return this.request('/api/chat', {
      method: 'POST',
      body: JSON.stringify(body),
    });
  },
  // Photo-to-Ticket: multipart upload, response sama dengan /api/chat
  async chatPhoto(file, sessionId, caption) {
    const fd = new FormData();
    fd.append('photo', file);
    fd.append('session_id', sessionId);
    const t = getToken();
    if (t) fd.append('token', t);
    if (caption) fd.append('message', caption);
    // Fallback token via query karena beberapa proxy/log memudahkan tracing
    const qs = t ? '?t=' + encodeURIComponent(t) : '';
    const res = await fetch('/api/chat/photo' + qs, {
      method: 'POST',
      headers: { ...guestHeaders() },
      body: fd,
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      const err = new Error(data?.error?.message || `Upload failed ${res.status}`);
      err.status = res.status;
      err.data = data;
      throw err;
    }
    return data;
  },
  createTicket(data) {
    const body = Object.assign({}, data);
    const t = getToken();
    if (t) body.token = t;
    return this.request('/api/tickets', {
      method: 'POST',
      body: JSON.stringify(body),
    });
  },
  getTickets(params = {}) {
    const p = Object.assign({}, params);
    const t = getToken();
    if (t && !p.t) p.t = t;
    const q = new URLSearchParams(p).toString();
    return this.request('/api/tickets' + (q ? '?' + q : ''));
  },
  getTicket(id) {
    const t = getToken();
    const q = t ? '?t=' + encodeURIComponent(t) : '';
    return this.request('/api/tickets/' + encodeURIComponent(id) + q);
  },
  // Identitas kamar tamu dari QR token
  getGuestMe() {
    return this.request('/api/guest/me');
  },
  hotelInfo(category) {
    const q = category ? '?category=' + encodeURIComponent(category) : '';
    return this.request('/api/hotel-info' + q);
  },
  recommendations(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.request('/api/recommendations' + (q ? '?' + q : ''));
  },
  // Room-service dining: menu catalog (public), place & list orders (guest).
  getMenu() {
    return this.request('/api/menu');
  },
  placeOrder(session_id, room_number, items, note) {
    const body = { session_id, room_number, items, note };
    const t = getToken();
    if (t) body.token = t;
    return this.request('/api/orders', {
      method: 'POST',
      body: JSON.stringify(body),
    });
  },
  getOrders(session_id) {
    const t = getToken();
    const q = new URLSearchParams();
    if (session_id) q.set('session_id', session_id);
    if (t) q.set('t', t);
    const s = q.toString();
    return this.request('/api/orders' + (s ? '?' + s : ''));
  },
  health() {
    return this.request('/api/health');
  }
};

// ---------- Guest token (QR) ----------
const GUEST_TOKEN_KEY = 'batiqa_guest_token';

// Capture ?t= di URL: simpan token lalu bersihkan address bar
(function initGuestTokenFromURL(){
  try{
    const params = new URLSearchParams(location.search);
    const t = params.get('t');
    if(t){
      localStorage.setItem(GUEST_TOKEN_KEY, t);
      params.delete('t');
      const rest = params.toString();
      history.replaceState(null, '', location.pathname + (rest ? '?' + rest : '') + location.hash);
    }
  }catch(e){ /* jangan crash */ }
})();

function getToken(){
  try{ return localStorage.getItem(GUEST_TOKEN_KEY) || ''; }catch(e){ return ''; }
}

// Sisipkan token ke URL navigasi internal app (header tak tersedia untuk <a>)
function withToken(url){
  const t = getToken();
  if(!t || !url) return url;
  if(/^(tel:|mailto:|javascript:)/i.test(url)) return url;
  if(url.charAt(0) === '#') return url;
  let u;
  try{ u = new URL(url, location.origin); }catch(e){ return url; }
  if(u.origin !== location.origin) return url; // link eksternal: biarkan
  const p = u.pathname;
  const isAppPage = p === '/' || p === '/index.html' || p.indexOf('/guest/') === 0;
  if(!isAppPage) return url;
  if(!u.searchParams.get('t')) u.searchParams.set('t', t);
  return u.pathname + u.search + u.hash;
}

// Header auth utama untuk semua request tamu
function guestHeaders(){
  const t = getToken();
  return t ? { 'X-Guest-Token': t } : {};
}

// Tambal semua <a href> internal agar selalu membawa token saat DOM siap
(function patchNavLinks(){
  function apply(){
    document.querySelectorAll('a[href]').forEach(a=>{
      const current = a.getAttribute('href');
      const next = withToken(current);
      if(next && next !== current) a.setAttribute('href', next);
    });
  }
  if(document.readyState === 'loading'){
    document.addEventListener('DOMContentLoaded', apply);
  } else {
    apply();
  }
})();

// Helpers
function getRoom() {
  const params = new URLSearchParams(location.search);
  return params.get('room') || params.get('room_number') || localStorage.getItem('batiqa_room') || '';
}
function getSession() {
  let s = localStorage.getItem('batiqa_session');
  if (!s) {
    s = 'guest-' + Date.now() + '-' + Math.random().toString(36).slice(2,7);
    localStorage.setItem('batiqa_session', s);
  }
  return s;
}
function setRoom(room) {
  if (room) localStorage.setItem('batiqa_room', room);
}
