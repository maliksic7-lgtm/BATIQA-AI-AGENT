// BATIQA API Client - lightweight, no heavy libs
const API = {
  base: '',
  async request(path, opts = {}) {
    const res = await fetch(path, {
      headers: { 'Content-Type': 'application/json', ...opts.headers },
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
    return this.request('/api/chat', {
      method: 'POST',
      body: JSON.stringify({ session_id, room_number, message }),
    });
  },
  createTicket(data) {
    return this.request('/api/tickets', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
  getTickets(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.request('/api/tickets' + (q ? '?' + q : ''));
  },
  getTicket(id) {
    return this.request('/api/tickets/' + encodeURIComponent(id));
  },
  hotelInfo(category) {
    const q = category ? '?category=' + encodeURIComponent(category) : '';
    return this.request('/api/hotel-info' + q);
  },
  recommendations(params = {}) {
    const q = new URLSearchParams(params).toString();
    return this.request('/api/recommendations' + (q ? '?' + q : ''));
  },
  health() {
    return this.request('/api/health');
  }
};

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
