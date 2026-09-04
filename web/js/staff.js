// Staff Dashboard - Desktop-first, operational efficiency, real backend data
(function(){
  const token = localStorage.getItem('batiqa_staff_token');
  const staff = JSON.parse(localStorage.getItem('batiqa_staff') || 'null');

  // Redirect to login if no token
  if(!token){
    location.href='/staff/login.html';
    return;
  }

  const staffInfo = document.getElementById('staffInfo');
  const deptBadge = document.getElementById('deptBadge');
  const tbody = document.getElementById('tbody');
  const tableWrap = document.getElementById('tableWrap');
  const loading = document.getElementById('loading');
  const empty = document.getElementById('empty');
  const errorEl = document.getElementById('error');
  const stats = {
    open: document.getElementById('statOpen'),
    high: document.getElementById('statHigh'),
    hk: document.getElementById('statHK'),
    eng: document.getElementById('statEng'),
    resolved: document.getElementById('statResolved'),
  };
  const deptFilter = document.getElementById('deptFilter');
  const statusFilter = document.getElementById('statusFilter');
  const priorityFilter = document.getElementById('priorityFilter');
  const modal = document.getElementById('modal');
  const modalClose = document.getElementById('modalClose');
  const modalCancel = document.getElementById('modalCancel');
  const modalSave = document.getElementById('modalSave');
  const modalError = document.getElementById('modalError');
  const assignBtn = document.getElementById('assignBtn');
  const mAssign = document.getElementById('mAssign');
  const liveDot = document.getElementById('liveDot');
  const qrModal = document.getElementById('qrModal');
  const qrClose = document.getElementById('qrClose');
  const qrBtn = document.getElementById('qrBtn');
  const qrImg = document.getElementById('qrImg');
  const qrRoomLabel = document.getElementById('qrRoomLabel');
  const toastsRoot = document.getElementById('toasts');

  let currentTicket = null;
  let ticketsCache = [];
  let myStaffId = null;
  let es = null;
  let qrObjectUrl = null;

  function authHeaders(){
    return { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' };
  }

  async function checkAuth(){
    try{
      const res = await fetch('/api/staff/me', {headers: authHeaders()});
      if(!res.ok) throw new Error('Unauthorized');
      const data = await res.json();
      staffInfo.textContent = data.name + ' • ' + data.department;
      deptBadge.textContent = data.department;
      myStaffId = data.id;
      // Pre-select dept filter to staff department for operational focus, but allow all
      // For ADMIN, show all
      if(data.department !== 'ADMIN' && !deptFilter.value){
        // Optional: filter to own department by default for HK/ENG
        // deptFilter.value = data.department;
      }
    }catch(e){
      localStorage.removeItem('batiqa_staff_token');
      localStorage.removeItem('batiqa_staff');
      location.href='/staff/login.html';
    }
  }

  async function loadStats(){
    try{
      const res = await fetch('/api/tickets/stats', {headers: authHeaders()});
      if(!res.ok) throw new Error('Failed stats');
      const data = await res.json();
      stats.open.textContent = data.total_open;
      stats.high.textContent = data.high_priority;
      stats.hk.textContent = data.housekeeping;
      stats.eng.textContent = data.engineering;
      stats.resolved.textContent = data.resolved_today;
    }catch(e){
      // Stats fail shouldn't block tickets
      console.error(e);
      stats.open.textContent='—';
    }
  }

  function showLoading(v){ loading.style.display = v?'flex':'none'; }
  function showError(msg){
    errorEl.textContent = msg;
    errorEl.style.display = 'block';
  }
  function hideError(){ errorEl.style.display='none'; }

  function esc(s){
    const d=document.createElement('div'); d.textContent=s; return d.innerHTML;
  }
  function fmtDate(iso){
    try{ const d=new Date(iso); return d.toLocaleString([], {month:'short', day:'numeric', hour:'2-digit', minute:'2-digit'}); }catch(e){ return iso; }
  }

  async function loadTickets(quiet){
    hideError();
    if(!quiet){
      showLoading(true);
      tableWrap.style.display='none';
    }
    empty.style.display='none';
    const params = new URLSearchParams();
    if(deptFilter.value) params.set('department', deptFilter.value);
    if(statusFilter.value) params.set('status', statusFilter.value);
    if(priorityFilter.value) params.set('priority', priorityFilter.value);
    // No room filter for staff - see all

    try{
      const res = await fetch('/api/tickets' + (params.toString()?'?'+params.toString():''), {headers: authHeaders()});
      // For List, auth not strictly required, but we send it
      // If 401, redirect to login
      if(res.status===401){
        localStorage.removeItem('batiqa_staff_token');
        location.href='/staff/login.html';
        return;
      }
      if(!res.ok){
        const err = await res.json().catch(()=>({error:{message:'Failed'}}));
        throw new Error(err?.error?.message||'Failed to load tickets');
      }
      const data = await res.json();
      const tickets = data.tickets || [];
      ticketsCache = tickets;
      showLoading(false);
      if(tickets.length===0){
        tableWrap.style.display='none';
        empty.style.display='block';
        return;
      }
      tbody.innerHTML='';
      tickets.forEach(t=>{
        const tr=document.createElement('tr');
        tr.innerHTML=`
          <td><strong style="color:var(--navy)">${esc(t.ticket_number)}</strong></td>
          <td>${esc(t.room_number)}</td>
          <td><span style="font-size:11px; font-weight:600; letter-spacing:0.06em;">${esc(t.department)}</span></td>
          <td>${esc(t.category)}</td>
          <td style="max-width:220px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;" title="${esc(t.description)}">${esc(t.description)}</td>
          <td><span class="priority priority--${esc(t.priority)}">${esc(t.priority)}</span></td>
          <td><span class="status status--${esc(t.status)}">${esc(t.status)}</span></td>
          <td style="white-space:nowrap;">${fmtDate(t.created_at)}</td>
          <td><button class="filter" data-id="${esc(t.ticket_number)}" aria-label="View ${esc(t.ticket_number)}">View</button></td>
        `;
        tbody.appendChild(tr);
      });
      tableWrap.style.display='block';
      // Attach view handlers
      tbody.querySelectorAll('button').forEach(btn=>{
        btn.addEventListener('click', ()=> openModal(btn.dataset.id));
      });
    }catch(e){
      showLoading(false);
      showError(e.message);
    }
  }

  function openModal(ticketNumber){
    const t = ticketsCache.find(x=>x.ticket_number===ticketNumber);
    if(!t) return;
    currentTicket = t;
    document.getElementById('mTicket').textContent = t.ticket_number;
    document.getElementById('mRoom').textContent = t.room_number;
    document.getElementById('mDept').textContent = t.department;
    document.getElementById('mCat').textContent = t.category;
    document.getElementById('mDesc').textContent = t.description;
    document.getElementById('mPriority').value = t.priority;
    document.getElementById('mCreated').textContent = fmtDate(t.created_at);
    document.getElementById('mStatus').value = t.status;
    mAssign.textContent = '—';
    loadAssignments(t.ticket_number);
    modalError.style.display='none';
    modal.classList.add('open');
  }
  function closeModal(){ modal.classList.remove('open'); currentTicket=null; }

  async function loadAssignments(ticketNumber){
    try{
      const res = await fetch('/api/tickets/' + encodeURIComponent(ticketNumber) + '/assignments');
      if(!res.ok) throw new Error();
      const data = await res.json();
      const list = data.assignments || [];
      mAssign.textContent = list.length ? ('Staff #' + list.map(a=>a.staff_id).join(', #')) : 'Unassigned';
    }catch(e){
      mAssign.textContent = '—';
    }
  }

  async function assignToMe(){
    if(!currentTicket) return;
    if(!myStaffId){
      modalError.textContent = 'Staff ID not available, please re-login.';
      modalError.style.display = 'block';
      return;
    }
    assignBtn.disabled = true;
    modalError.style.display='none';
    try{
      const res = await fetch('/api/tickets/' + encodeURIComponent(currentTicket.ticket_number) + '/assign', {
        method:'POST',
        headers: authHeaders(),
        body: JSON.stringify({staff_id: myStaffId})
      });
      const data = await res.json().catch(()=>({}));
      if(!res.ok) throw new Error(data?.error?.message||'Failed to assign');
      await loadAssignments(currentTicket.ticket_number);
    }catch(e){
      modalError.textContent = e.message;
      modalError.style.display = 'block';
    }finally{
      assignBtn.disabled = false;
    }
  }

  async function saveTicket(){
    if(!currentTicket) return;
    const newStatus = document.getElementById('mStatus').value;
    const newPriority = document.getElementById('mPriority').value;
    const statusChanged = newStatus !== currentTicket.status;
    const priorityChanged = newPriority !== currentTicket.priority;
    if(!statusChanged && !priorityChanged){
      closeModal(); return;
    }
    modalSave.disabled=true;
    modalSave.textContent='Updating...';
    modalError.style.display='none';
    try{
      // Sequential: status first (transition validation), then priority
      if(statusChanged){
        const res = await fetch('/api/tickets/' + encodeURIComponent(currentTicket.ticket_number) + '/status', {
          method:'PATCH',
          headers: authHeaders(),
          body: JSON.stringify({status: newStatus})
        });
        const data = await res.json().catch(()=>({}));
        if(!res.ok) throw new Error(data?.error?.message||'Failed to update status');
      }
      if(priorityChanged){
        const res2 = await fetch('/api/tickets/' + encodeURIComponent(currentTicket.ticket_number) + '/priority', {
          method:'PATCH',
          headers: authHeaders(),
          body: JSON.stringify({priority: newPriority})
        });
        const data2 = await res2.json().catch(()=>({}));
        if(!res2.ok) throw new Error(data2?.error?.message||'Failed to update priority');
      }
      // Success - close and reload
      closeModal();
      await loadTickets();
      await loadStats();
    }catch(e){
      modalError.textContent = e.message;
      modalError.style.display='block';
    }finally{
      modalSave.disabled=false;
      modalSave.textContent='Update Ticket';
    }
  }

  // ---------- Toasts ----------
  function showToast(msg){
    if(!toastsRoot) return;
    while(toastsRoot.children.length >= 3){
      toastsRoot.removeChild(toastsRoot.firstChild);
    }
    const el = document.createElement('div');
    el.className = 'toast';
    el.setAttribute('role', 'status');
    el.textContent = msg;
    let gone = false;
    const dismiss = () => {
      if(gone) return;
      gone = true;
      clearTimeout(timer);
      el.classList.add('hide');
      setTimeout(()=>{ el.remove(); }, 320);
    };
    const timer = setTimeout(dismiss, 5000);
    el.addEventListener('click', dismiss);
    toastsRoot.appendChild(el);
  }

  // ---------- Live feed (SSE) ----------
  function setLive(on){
    if(!liveDot) return;
    liveDot.classList.toggle('is-live', !!on);
    const txt = liveDot.querySelector('.live-status__text');
    if(txt) txt.textContent = on ? 'Live' : 'Offline';
    liveDot.title = on
      ? 'Realtime terhubung'
      : 'Koneksi realtime terputus — menyambung ulang otomatis';
  }

  function initEvents(){
    if(es || typeof EventSource === 'undefined') return;
    try{
      es = new EventSource('/api/events?t=' + encodeURIComponent(token));
    }catch(e){
      console.error('EventSource gagal dibuat', e);
      return;
    }
    es.onopen = () => setLive(true);
    es.onerror = () => setLive(false); // browser auto-reconnect
    es.onmessage = (ev) => {
      let msg = null;
      try{ msg = JSON.parse(ev.data); }catch(_){ return; }
      if(!msg || !msg.type) return;
      if(msg.type === 'ticket.created'){
        const t = msg.ticket || {};
        showToast('Tiket baru ' + (t.ticket_number || '?') + ' • Kamar ' + (t.room_number || '-'));
        loadStats();
        loadTickets(true);
        loadAnalytics(true);
        loadInfographics(true);
      }else if(msg.type === 'ticket.updated'){
        loadStats();
        loadTickets(true);
        loadAnalytics(true);
        loadInfographics(true);
      }
    };
  }

  // ---------- QR Kamar ----------
  async function openQr(){
    const room = (prompt('Nomor kamar (contoh: 305):') || '').trim();
    if(!room) return;
    qrRoomLabel.textContent = 'Kamar ' + room;
    qrImg.removeAttribute('src');
    qrModal.classList.add('open');
    try{
      const res = await fetch('/api/rooms/' + encodeURIComponent(room) + '/qr?t=' + encodeURIComponent(token), {
        headers: { 'Authorization': 'Bearer ' + token }
      });
      if(!res.ok) throw new Error('Gagal memuat QR kamar ' + room + ' (' + res.status + ')');
      const blob = await res.blob();
      if(qrObjectUrl){ URL.revokeObjectURL(qrObjectUrl); qrObjectUrl = null; }
      qrObjectUrl = URL.createObjectURL(blob);
      qrImg.src = qrObjectUrl;
    }catch(e){
      closeQr();
      showError(e.message);
    }
  }

  function closeQr(){
    qrModal.classList.remove('open');
    if(qrObjectUrl){
      URL.revokeObjectURL(qrObjectUrl);
      qrObjectUrl = null;
    }
  }

  // ---------- Analytics ----------
  async function loadAnalytics(quiet){
    try{
      const res = await fetch('/api/analytics', {headers: authHeaders()});
      if(!res.ok) throw new Error('Failed analytics');
      buildCharts(await res.json());
    }catch(e){
      if(!quiet) console.error(e);
    }
  }

  // ---------- Infographics ----------
  async function loadInfographics(quiet){
    try{
      const res = await fetch('/api/analytics/infographics', {headers: authHeaders()});
      if(!res.ok) throw new Error('Failed infographics');
      buildInfographics(await res.json());
    }catch(e){
      if(!quiet) console.error(e);
    }
  }

  function buildInfographics(data){
    const d = data || {};
    renderCategoryList('topComplaints', d.top_complaints, 'hbar-fill--high');
    renderCategoryList('topBorrowed', d.top_borrowed, 'hbar-fill--low');
    renderCategoryList('topAsked', d.top_asked, 'hbar-fill--gold');
    renderCategoryList('topOrdered', d.top_ordered, '');
  }

  function renderCategoryList(id, list, fillClass){
    const el = document.getElementById(id);
    if(!el) return;
    const arr = (list || []).slice(0, 5);
    if(!arr.length){
      el.innerHTML = '<p class="an-empty">Belum ada data</p>';
      return;
    }
    const max = Math.max(1, ...arr.map(c => Number(c.count) || 0));
    el.innerHTML = arr.map((c, i) => {
      const count = Number(c.count) || 0;
      const pct = Math.round((count / max) * 100);
      return '<div class="cat-row">'
        + '<span class="cat-rank">' + (i + 1) + '</span>'
        + '<div class="cat-main">'
        + '<div class="hbar-top"><span>' + esc(labelize(c.category)) + '</span><strong>' + count + '</strong></div>'
        + '<div class="hbar-track"><div class="hbar-fill ' + (fillClass || '') + '" style="width:' + pct + '%"></div></div>'
        + '</div>'
        + '</div>';
    }).join('');
  }

  function labelize(s){
    return String(s || '').split('_')
      .map(w => w ? w.charAt(0) + w.slice(1).toLowerCase() : w)
      .join(' ');
  }

  function shortDate(iso){
    const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(String(iso || ''));
    if(m) return String(Number(m[3])) + '/' + String(Number(m[2]));
    try{
      const d = new Date(iso);
      return d.getDate() + '/' + (d.getMonth() + 1);
    }catch(e){ return iso; }
  }

  function niceStep(maxV){
    const raw = Math.max(maxV / 4, 1);
    const pow = Math.pow(10, Math.floor(Math.log10(raw)));
    for(const m of [1, 2, 2.5, 5, 10]){
      if(m * pow >= raw) return m * pow;
    }
    return 10 * pow;
  }

  function buildBarChartSvg(daily){
    const days = (daily || []).map(d => ({
      label: shortDate(d.date),
      created: Number(d.created) || 0,
      resolved: Number(d.resolved) || 0,
    }));
    if(!days.length) return '<p class="an-empty">Belum ada data</p>';
    const W = 580, H = 232, L = 34, R = 8, T = 10, B = 30;
    const pw = W - L - R, ph = H - T - B;
    const maxV = Math.max(1, ...days.map(d => Math.max(d.created, d.resolved)));
    const step = niceStep(maxV);
    const y = v => T + ph - (v / maxV) * ph;
    const slot = pw / days.length;
    const bw = Math.max(8, Math.min(16, slot * 0.28));
    let s = '<svg class="barchart" viewBox="0 0 ' + W + ' ' + H + '" role="img" aria-label="Tiket dibuat vs selesai, ' + days.length + ' hari terakhir">';
    for(let v = 0; v <= maxV + step * 0.01; v += step){
      const yy = y(v);
      s += '<line x1="' + L + '" x2="' + (W - R) + '" y1="' + yy + '" y2="' + yy + '" stroke="rgba(3,60,90,' + (v === 0 ? 0.25 : 0.08) + ')" stroke-width="1"/>'
        + '<text x="' + (L - 7) + '" y="' + (yy + 3) + '" text-anchor="end" font-size="9" fill="#6E6A61">' + Math.round(v) + '</text>';
    }
    days.forEach((d, i) => {
      const cx = L + slot * i + slot / 2;
      const base = T + ph;
      const hc = y(d.created), hr = y(d.resolved);
      s += '<g>'
        + '<rect x="' + (cx - bw - 2) + '" y="' + hc + '" width="' + bw + '" height="' + Math.max(1, base - hc) + '" rx="3" fill="#033C5A"><title>' + d.label + ' · Created: ' + d.created + '</title></rect>'
        + '<rect x="' + (cx + 2) + '" y="' + hr + '" width="' + bw + '" height="' + Math.max(1, base - hr) + '" rx="3" fill="#C9AC85"><title>' + d.label + ' · Resolved: ' + d.resolved + '</title></rect>'
        + '</g>';
      s += '<text x="' + cx + '" y="' + (H - 9) + '" text-anchor="middle" font-size="9.5" fill="#6E6A61">' + esc(d.label) + '</text>';
    });
    s += '</svg>';
    return s;
  }

  function hbarColorClass(key, priorityMode){
    const k = String(key || '').toUpperCase();
    if(priorityMode){
      if(k === 'HIGH') return 'hbar-fill--high';
      if(k === 'LOW') return 'hbar-fill--low';
      return 'hbar-fill--gold';
    }
    if(k === 'ENGINEERING') return 'hbar-fill--gold';
    if(k === 'FRONT_OFFICE') return 'hbar-fill--low';
    return '';
  }

  function renderHBars(id, obj, priorityMode){
    const el = document.getElementById(id);
    if(!el) return;
    const entries = Object.entries(obj || {}).sort((a, b) => b[1] - a[1]);
    if(!entries.length){
      el.innerHTML = '<p class="an-empty">Tidak ada tiket aktif</p>';
      return;
    }
    const max = Math.max(1, ...entries.map(e => e[1]));
    el.innerHTML = entries.map(([k, v]) => {
      const pct = Math.round((v / max) * 100);
      return '<div class="hbar-row">'
        + '<div class="hbar-top"><span>' + esc(labelize(k)) + '</span><strong>' + esc(String(v)) + '</strong></div>'
        + '<div class="hbar-track"><div class="hbar-fill ' + hbarColorClass(k, priorityMode) + '" style="width:' + pct + '%"></div></div>'
        + '</div>';
    }).join('');
  }

  function renderCategories(id, list){
    const el = document.getElementById(id);
    if(!el) return;
    const arr = (list || []).slice(0, 5);
    if(!arr.length){
      el.innerHTML = '<p class="an-empty">Belum ada tiket</p>';
      return;
    }
    const max = Math.max(1, ...arr.map(c => Number(c.count) || 0));
    el.innerHTML = arr.map((c, i) => {
      const count = Number(c.count) || 0;
      const pct = Math.round((count / max) * 100);
      return '<div class="cat-row">'
        + '<span class="cat-rank">' + (i + 1) + '</span>'
        + '<div class="cat-main">'
        + '<div class="hbar-top"><span>' + esc(labelize(c.category)) + '</span><strong>' + count + '</strong></div>'
        + '<div class="hbar-track"><div class="hbar-fill hbar-fill--gold" style="width:' + pct + '%"></div></div>'
        + '</div>'
        + '</div>';
    }).join('');
  }

  function buildCharts(data){
    const d = data || {};
    const avgEl = document.getElementById('avgResolution');
    if(avgEl) avgEl.textContent = (typeof d.avg_resolution_hours === 'number' ? d.avg_resolution_hours.toFixed(1) : '—') + 'h';
    const openEl = document.getElementById('anTotalOpen');
    if(openEl && typeof d.total_open !== 'undefined') openEl.textContent = d.total_open;
    const totalEl = document.getElementById('anTotalTickets');
    if(totalEl && typeof d.total_tickets !== 'undefined') totalEl.textContent = d.total_tickets;
    const chartEl = document.getElementById('barChart');
    if(chartEl) chartEl.innerHTML = buildBarChartSvg(d.daily);
    renderHBars('deptBars', d.active_by_department, false);
    renderHBars('priorityBars', d.active_by_priority, true);
    renderCategories('topCats', d.top_categories);
  }

  // Events
  document.getElementById('refreshBtn').addEventListener('click', ()=>{ loadTickets(); loadStats(); loadAnalytics(); loadInfographics(); });
  document.getElementById('clearFilter').addEventListener('click', ()=>{
    deptFilter.value=''; statusFilter.value=''; priorityFilter.value='';
    loadTickets();
  });
  deptFilter.addEventListener('change', loadTickets);
  statusFilter.addEventListener('change', loadTickets);
  priorityFilter.addEventListener('change', loadTickets);
  modalClose.addEventListener('click', closeModal);
  modalCancel.addEventListener('click', closeModal);
  modalSave.addEventListener('click', saveTicket);
  assignBtn.addEventListener('click', assignToMe);
  modal.addEventListener('click', (e)=>{ if(e.target===modal) closeModal(); });
  qrBtn.addEventListener('click', openQr);
  qrClose.addEventListener('click', closeQr);
  qrModal.addEventListener('click', (e)=>{ if(e.target===qrModal) closeQr(); });
  document.addEventListener('keydown', (e)=>{
    if(e.key==='Escape'){
      if(modal.classList.contains('open')) closeModal();
      if(qrModal.classList.contains('open')) closeQr();
    }
  });
  document.getElementById('logoutBtn').addEventListener('click', async (e)=>{
    e.preventDefault();
    try{ await fetch('/api/staff/logout', {method:'POST', headers: authHeaders()}); }catch(_){}
    localStorage.removeItem('batiqa_staff_token');
    localStorage.removeItem('batiqa_staff');
    location.href='/staff/login.html';
  });

  // Init
  checkAuth().then(()=>{
    loadStats();
    loadTickets();
    loadAnalytics();
    loadInfographics();
    initEvents();
  });
  // Polling fallback (SSE adalah mekanisme utama realtime)
  setInterval(()=>{ loadStats(); loadTickets(); }, 60000);
})();
