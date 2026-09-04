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

  // Photo upload elements
  const photoBtn = document.getElementById('photoBtn');
  const photoInput = document.getElementById('photoInput');
  const photoPreview = document.getElementById('photoPreview');
  const photoPreviewImg = document.getElementById('photoPreviewImg');
  const photoPreviewName = document.getElementById('photoPreviewName');
  const photoRemove = document.getElementById('photoRemove');

  const MAX_PHOTO_SIZE = 8 * 1024 * 1024; // 8MB
  let pendingFile = null;
  let pendingThumbURL = null;

  let room = getRoom();
  const session = getSession();
  if(room) roomLabel.textContent = 'Room ' + room;

  // Identitas kamar via guest token QR (fallback bila ?room tidak ada)
  if(!room && typeof API.getGuestMe === 'function'){
    API.getGuestMe()
      .then(me => {
        if(me && me.room_number){
          room = me.room_number;
          setRoom(room);
          roomLabel.textContent = 'Room ' + room;
        }
      })
      .catch(()=>{ /* token invalid / offline - chat tetap jalan */ });
  }

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
    if(photoBtn) photoBtn.disabled = v;
  }
  function showError(msg){
    errorText.textContent = msg;
    errorBar.style.display = 'flex';
  }
  function hideError(){ errorBar.style.display='none'; }

  // Pesan ramah untuk sesi/token invalid - jangan crash
  function friendlyError(e){
    if(e && (e.status === 401 || e.status === 403)){
      return 'Sesi tidak valid. Silakan scan ulang QR di kamar Anda.';
    }
    return (e && e.message) || 'Terjadi kesalahan. Silakan coba lagi.';
  }

  // Response {message,intent,requires_ticket,ticket_id} -> bubble AI + kartu tiket
  async function handleChatResponse(res){
    let ticket = null;
    if(res.ticket_id){
      try{
        ticket = await API.getTicket(res.ticket_id);
      }catch(e){
        ticket = { ticket_number: res.ticket_id, room_number: room, department: '', priority: '' };
      }
    }
    addBubble(res.message, 'ai', ticket, res.intent);
  }

  async function sendMessage(text){
    const msg = (text || input.value).trim();
    if(pendingFile){
      await sendPhoto(msg);
      return;
    }
    if(!msg) return;
    lastMessage = msg;
    hideError();
    addBubble(msg,'user');
    input.value = '';
    input.style.height='42px';
    setLoading(true);

    try{
      const res = await API.chat(session, room, msg);
      await handleChatResponse(res);
    }catch(e){
      // Fallback per AI failure - chat not broken
      console.error(e);
      showError(friendlyError(e));
      // Still show fallback bubble so chat not broken
      const fallback = room ? 'Maaf, layanan AI sedang mengalami gangguan. Silakan coba lagi.' : 'Maaf, layanan AI sedang gangguan. Boleh saya tahu nomor kamar Anda?';
      addBubble(fallback, 'ai');
    }finally{
      setLoading(false);
      input.focus();
    }
  }

  // ---------- Photo-to-Ticket ----------
  function clearPendingPhoto(){
    pendingFile = null;
    if(pendingThumbURL){ URL.revokeObjectURL(pendingThumbURL); pendingThumbURL = null; }
    photoInput.value = '';
    if(photoPreview) photoPreview.style.display = 'none';
    if(photoPreviewImg) photoPreviewImg.removeAttribute('src');
  }

  if(photoBtn && photoInput){
    photoBtn.addEventListener('click', ()=>{ if(!pendingFile) photoInput.click(); });
    photoInput.addEventListener('change', ()=>{
      const f = photoInput.files && photoInput.files[0];
      if(!f) return;
      if(f.size > MAX_PHOTO_SIZE){
        showError('Ukuran foto maksimal 8MB.');
        clearPendingPhoto();
        return;
      }
      if(f.type && !/^image\/(jpeg|jpg|png|webp)$/i.test(f.type)){
        showError('Format foto harus JPG, PNG, atau WebP.');
        clearPendingPhoto();
        return;
      }
      pendingFile = f;
      if(photoPreviewImg){
        pendingThumbURL = URL.createObjectURL(f);
        photoPreviewImg.src = pendingThumbURL;
      }
      if(photoPreviewName) photoPreviewName.textContent = f.name || 'Foto kerusakan';
      if(photoPreview) photoPreview.style.display = 'flex';
      hideError();
      input.focus();
    });
    if(photoRemove){
      photoRemove.addEventListener('click', ()=>{
        clearPendingPhoto();
        input.focus();
      });
    }
  }

  async function sendPhoto(caption){
    if(!pendingFile) return;
    lastMessage = caption || '';
    hideError();
    addBubble('[\u{1F4F7} Foto dikirim]' + (caption ? '\n' + caption : ''), 'user');
    setLoading(true);
    const file = pendingFile;
    try{
      const res = await API.chatPhoto(file, session, caption || '');
      await handleChatResponse(res);
      clearPendingPhoto();
    }catch(e){
      console.error(e);
      showError(friendlyError(e));
      // Preview dipertahankan agar tamu bisa kirim ulang
      const fallback = room ? 'Maaf, foto gagal diproses. Silakan coba lagi.' : 'Maaf, layanan AI sedang gangguan. Boleh saya tahu nomor kamar Anda?';
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

  // Restore chat history from server (guest may reload the page or re-scan QR)
  (async function restoreHistory(){
    if(!session) return;
    try{
      const t = getToken();
      const url = '/api/conversations?session_id=' + encodeURIComponent(session) + '&limit=50' + (t ? '&t=' + encodeURIComponent(t) : '');
      const res = await fetch(url);
      if(!res.ok) return;
      const data = await res.json();
      const msgs = data.messages || [];
      if(!msgs.length) return;
      messagesEl.innerHTML='';
      msgs.forEach(m=>{
        addBubble(m.message, m.role==='user' ? 'user' : 'ai', null, m.intent);
      });
    }catch(e){
      // History is best-effort; never block new chat
      console.warn('history unavailable');
    }
  })();
})();
