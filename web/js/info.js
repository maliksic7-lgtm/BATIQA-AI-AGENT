// Hotel Info + Recommendations - handles loading/error/empty, fast, no heavy libs
(function(){
  const list = document.getElementById('list');
  const loading = document.getElementById('loading');
  const errorEl = document.getElementById('error');
  const errorText = document.getElementById('errorText');
  const empty = document.getElementById('empty');
  const filter = document.getElementById('filter');
  const badge = document.getElementById('roomBadge');
  const recList = document.getElementById('recList');
  const recLoading = document.getElementById('recLoading');

  const room = getRoom();
  if(room){ badge.textContent='Room '+room; badge.style.display='block'; }

  let currentCat = '';

  function esc(s){ const d=document.createElement('div'); d.textContent=s; return d.innerHTML; }

  async function loadInfo(){
    loading.style.display='flex';
    errorEl.style.display='none';
    empty.style.display='none';
    list.style.display='none';
    try{
      const data = await API.hotelInfo(currentCat||undefined);
      const items = data.items||[];
      loading.style.display='none';
      if(items.length===0){ empty.style.display='block'; return; }
      list.innerHTML='';
      items.forEach(it=>{
        const div=document.createElement('div');
        div.className='info-card';
        div.innerHTML=`<p class="info-card__cat">${esc(it.category)}</p><p class="info-card__title">${esc(it.title)}</p><p class="info-card__content">${esc(it.content)}</p>`;
        list.appendChild(div);
      });
      list.style.display='grid';
    }catch(e){
      loading.style.display='none';
      errorText.textContent=e.message;
      errorEl.style.display='flex';
    }
  }

  async function loadRec(){
    recLoading.style.display='flex';
    recList.style.display='none';
    try{
      const data = await API.recommendations({});
      const items = data.items||[];
      recLoading.style.display='none';
      if(items.length===0) return;
      recList.innerHTML='';
      items.slice(0,4).forEach(it=>{
        const div=document.createElement('div');
        div.className='info-card';
        const price = it.price_min!=null ? `Rp${it.price_min} - ${it.price_max}` : '—';
        const dist = it.distance_km!=null ? `${it.distance_km} km` : '';
        div.innerHTML=`<p class="info-card__cat">${esc(it.category)}</p><p class="info-card__title">${esc(it.name)}</p><p class="info-card__content">${price} ${dist? '• '+dist : ''}</p>`;
        recList.appendChild(div);
      });
      recList.style.display='grid';
    }catch(e){
      recLoading.style.display='none';
      // silent fail for rec
    }
  }

  filter.addEventListener('click', (e)=>{
    if(e.target.classList.contains('chip')){
      filter.querySelectorAll('.chip').forEach(c=>c.classList.remove('active'));
      e.target.classList.add('active');
      currentCat=e.target.dataset.cat;
      loadInfo();
    }
  });

  loadInfo();
  loadRec();
})();
