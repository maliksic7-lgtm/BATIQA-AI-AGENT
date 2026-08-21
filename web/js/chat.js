// Chat logic - mobile-first, robust to AI failure
(function(){
  const messagesEl = document.getElementById('messages');
  const input = document.getElementById('msgInput');
  const form = document.getElementById('chatForm');
  const loading = document.getElementById('loading');
  const errorBar = document.getElementById('errorBar');
  const errorText = document.getElementById('errorText');
  const retryBtn = document.getElementById('retryBtn');
  const suggestions = document.getElementById('suggestions');
  const roomLabel = document.getElementById('roomLabel');

  const room = getRoom();
  const session = getSession();
  if(room) roomLabel.textContent = 'Room ' + room;

  // Quick param ?quick=towel
  const params = new URLSearchParams(location.search);
  const quick = params.get('quick');
  if(quick){
    if(quick==='towel') input.value = 'Tolong antar 2 handuk ke kamar ' + (room||'305');
    if(quick==='ac') input.value = 'AC kamar saya tidak dingin';
    if(quick==='wifi') input.value = 'WiFi tidak bekerja';
  }
  // Also support ?room
  if(params.get('room')) setRoom(params.get('room'));

  let lastMessage = '';

  function addBubble(text, who, ticket, intent){
    const div = document.createElement('div');
    div.className = 'bubble ' + (who==='user' ? 'bubble--user' : 'bubble--ai');
    // Escape HTML
    const esc = document.createElement('div');
    esc.textContent = text;
    div.innerHTML = esc.innerHTML.replace(/\n/g,'<br>');

    if(who==='ai'){
      const meta = document.createElement('div');
      meta.className='bubble__meta';
      meta.textContent = 'BATIQA AI • ' + new Date().toLocaleTimeString([], {hour:'2-digit',minute:'2-digit'});
      div.appendChild(meta);
      if(ticket){
        const box = document.createElement('div');
        box.className='bubble__ticket';
        box.innerHTML = `<strong>Ticket ${ticket.ticket_number || ticket}</strong><div class="bubble__ticket__row"><span>Room ${ticket.room_number||room}</span><span class="priority priority--${ticket.priority||'MEDIUM'}">${ticket.priority||''}</span></div><div style="margin-top:4px;">${ticket.department||''} • ${ticket.category||intent||''}</div>`;
        div.appendChild(box);
      }
      if(intent && ticket){
        // also show priority
      }
    } else {
      // user meta not needed
    }
    messagesEl.appendChild(div);
    messagesEl.scrollTop = messagesEl.scrollHeight;
  }

  function setLoading(v){
    loading.style.display = v ? 'flex' : 'none';
    document.getElementById('sendBtn').disabled = v;
  }
  function showError(msg){
    errorText.textContent = msg;
    errorBar.style.display = 'flex';
  }
  function hideError(){ errorBar.style.display='none'; }

  async function sendMessage(text){
    const msg = (text || input.value).trim();
    if(!msg) return;
    lastMessage = msg;
    hideError();
    addBubble(msg,'user');
    input.value = '';
    input.style.height='42px';
    setLoading(true);

    try{
      const res = await API.chat(session, room, msg);
      // res: {message, intent, requires_ticket, ticket_id}
      let ticket = null;
      if(res.ticket_id){
        // Fetch ticket detail for display
        try{
          const detail = await API.getTicket(res.ticket_id);
          ticket = detail;
        }catch(e){
          ticket = { ticket_number: res.ticket_id, room_number: room, department: '', priority: '' };
        }
      } else if(res.requires_ticket && !res.ticket_id){
        // Missing room case - AI asked for room
        ticket = null;
      }
      addBubble(res.message, 'ai', ticket, res.intent);
    }catch(e){
      // Fallback per AI failure - chat not broken
      console.error(e);
      showError(e.message || 'Failed to send. Please try again.');
      // Still show fallback bubble so chat not broken
      const fallback = room ? 'Maaf, layanan AI sedang mengalami gangguan. Silakan coba lagi.' : 'Maaf, layanan AI sedang gangguan. Boleh saya tahu nomor kamar Anda?';
      addBubble(fallback, 'ai');
    }finally{
      setLoading(false);
      input.focus();
    }
  }

  // Auto-resize textarea
  input.addEventListener('input', ()=>{
    input.style.height='auto';
    input.style.height = Math.min(input.scrollHeight, 100) + 'px';
  });
  input.addEventListener('keydown', (e)=>{
    if(e.key==='Enter' && !e.shiftKey){
      e.preventDefault();
      sendMessage();
    }
  });
  form.addEventListener('submit', (e)=>{
    e.preventDefault();
    sendMessage();
  });
  suggestions.addEventListener('click', (e)=>{
    if(e.target.classList.contains('chip')){
      sendMessage(e.target.dataset.msg);
    }
  });
  retryBtn.addEventListener('click', ()=>{
    if(lastMessage) sendMessage(lastMessage);
  });

  // Accessibility: focus input on load (mobile will show keyboard, so delay)
  setTimeout(()=>{ if(window.innerWidth>480) input.focus(); }, 300);
})();
