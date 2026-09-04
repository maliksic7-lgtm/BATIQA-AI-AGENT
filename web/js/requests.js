// My Requests - mobile-first, handles loading/error/empty
(function(){
  const listEl = document.getElementById('list');
  const loading = document.getElementById('loading');
  const errorEl = document.getElementById('error');
  const errorText = document.getElementById('errorText');
  const empty = document.getElementById('empty');
  const badge = document.getElementById('roomBadge');
  const refreshBtn = document.getElementById('refreshBtn');

  let room = getRoom();
  function setBadge(){
    if(room){
      badge.textContent='Room '+room;
      badge.style.display='block';
    }
  }
  setBadge();

  // Identitas kamar via guest token QR (fallback bila ?room tidak ada)
  if(!room && typeof API.getGuestMe === 'function'){
    API.getGuestMe()
      .then(me=>{
        if(me && me.room_number){
          room = me.room_number;
          setRoom(room);
          setBadge();
          loadTickets();
        }
      })
      .catch(()=>{ /* token invalid / offline */ });
  }

  function showLoading(v){ loading.style.display = v?'flex':'none'; }
  function showError(msg){
    errorText.textContent = msg;
    errorEl.style.display='flex';
  }
  function hideError(){ errorEl.style.display='none'; }

  function statusClass(s){
    if(s==='OPEN') return '';
    if(s==='IN_PROGRESS') return 'status--active';
    if(s==='RESOLVED') return 'status--done';
    return '';
  }
  function priorityClass(p){
    return 'priority--'+ (p||'MEDIUM');
  }
  function formatDate(iso){
    try{
      const d=new Date(iso);
      return d.toLocaleString([], {month:'short', day:'numeric', hour:'2-digit', minute:'2-digit'});
    }catch(e){ return iso; }
  }

  async function loadTickets(){
    hideError();
    empty.style.display='none';
    listEl.style.display='none';
    showLoading(true);
    try{
      const params = {};
      if(room) params.room = room;
      // If no room, fetch all but will be empty for guest without room - show empty
      const data = await API.getTickets(params);
      const tickets = data.tickets || [];
      if(statusWatcher) statusWatcher.refresh(tickets);
      showLoading(false);
      if(tickets.length===0){
        empty.style.display='block';
        return;
      }
      listEl.innerHTML='';
      tickets.forEach(t=>{
        const div=document.createElement('div');
        div.className='ticket';
        div.innerHTML=`
          <div class="ticket__head">
            <span class="ticket__number">#${t.ticket_number}</span>
            <span class="ticket__dept">${t.department}</span>
          </div>
          <p class="ticket__desc">${escapeHtml(t.description)}</p>
          <div class="ticket__meta">
            <span>Room ${escapeHtml(t.room_number)}</span>
            <span class="priority ${priorityClass(t.priority)}">${t.priority}</span>
            <span>${formatDate(t.created_at)}</span>
          </div>
          <div class="ticket__status">
            <span class="status ${t.status==='OPEN'?'status--active':''}">Open</span>
            <span class="status ${t.status==='IN_PROGRESS'?'status--active':''}">In Progress</span>
            <span class="status ${t.status==='RESOLVED'?'status--done':''}">Resolved</span>
          </div>
        `;
        // Highlight active status
        // For visual, set active class based on actual status
        // Already handled: OPEN active is first, IN_PROGRESS second, RESOLVED third
        // If CANCELLED, show differently
        if(t.status==='CANCELLED'){
          div.querySelectorAll('.status').forEach(s=>s.classList.remove('status--active','status--done'));
        }
        listEl.appendChild(div);
      });
      listEl.style.display='grid';
    }catch(e){
      showLoading(false);
      const invalidSession = e && (e.status === 401 || e.status === 403);
      showError(invalidSession
        ? 'Sesi tidak valid. Silakan scan ulang QR di kamar Anda.'
        : (e.message || 'Failed to load. Check connection.'));
      console.error(e);
    }
  }

  function escapeHtml(s){
    const d=document.createElement('div');
    d.textContent=s;
    return d.innerHTML;
  }

  refreshBtn.addEventListener('click', loadTickets);

  // Real-time status notification watcher (uses foreground Web Notification when
  // a ticket moves to a new step). Best-effort; never blocks page.
  let statusWatcher = null;
  if(typeof window.BATIQA !== 'undefined' && window.BATIQA.watchTicketStatus){
    statusWatcher = window.BATIQA.watchTicketStatus(30000);
  }

  loadTickets();
  // Auto refresh every 15s for status updates ( Guest Notification flow )
  setInterval(loadTickets, 15000);
})();
