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

  let currentTicket = null;
  let ticketsCache = [];
  let myStaffId = null;

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

  async function loadTickets(){
    hideError();
    showLoading(true);
    tableWrap.style.display='none';
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

  // Events
  document.getElementById('refreshBtn').addEventListener('click', ()=>{ loadTickets(); loadStats(); });
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
  document.getElementById('logoutBtn').addEventListener('click', async (e)=>{
    e.preventDefault();
    try{ await fetch('/api/staff/logout', {method:'POST', headers: authHeaders()}); }catch(_){}
    localStorage.removeItem('batiqa_staff_token');
    localStorage.removeItem('batiqa_staff');
    location.href='/staff/login.html';
  });
  document.addEventListener('keydown', (e)=>{ if(e.key==='Escape' && modal.classList.contains('open')) closeModal(); });

  // Init
  checkAuth().then(()=>{ loadStats(); loadTickets(); });
  // Poll every 15s for operational efficiency (real-time-ish)
  setInterval(()=>{ loadStats(); loadTickets(); }, 15000);
})();
