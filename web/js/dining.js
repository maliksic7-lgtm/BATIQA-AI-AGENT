// Room dining: load menu, maintain cart, place order, list guest's orders.
(function () {
  const cart = {};
  let menu = [];

  const menuEl = document.getElementById('menu');
  const cartEl = document.getElementById('cartItems');
  const cartEmpty = document.getElementById('cartEmpty');
  const totalEl = document.getElementById('cartTotal');
  const statusEl = document.getElementById('orderStatus');
  const submitBtn = document.getElementById('placeOrderBtn');
  const badge = document.getElementById('roomBadge');
  const catRow = document.getElementById('catRow');
  const ordersEl = document.getElementById('orders');

  let room = getRoom();
  const session = getSession();

  function setBadge() {
    if (room) { badge.textContent = 'Room ' + room; badge.style.display = 'block'; }
  }
  setBadge();
  if (!room && typeof API.getGuestMe === 'function') {
    API.getGuestMe().then(function (me) {
      if (me && me.room_number) { room = me.room_number; setRoom(room); setBadge(); }
    }).catch(function () {});
  }

  // ---- catalog rendering ----
  function formatRupiah(n) {
    return 'Rp ' + (n || 0).toLocaleString('id-ID');
  }

  function renderMenu(filter) {
    const list = menu.filter(function (m) { return !filter || m.category === filter; });
    if (list.length === 0) { menuEl.innerHTML = '<p class="cart__empty">Tidak ada item untuk kategori ini.</p>'; return; }
    menuEl.innerHTML = '';
    list.forEach(function (m) {
      const card = document.createElement('div');
      card.className = 'dish';
      const inCart = cart[m.name] || 0;
      card.innerHTML =
        '<div class="dish__title"><span>' + escapeHtml(m.name) + '</span><span class="dish__tag">' + (m.category === 'DRINK' ? 'Minuman' : 'Makanan') + '</span></div>' +
        '<div class="dish__price">' + formatRupiah(m.price) + '</div>' +
        '<div class="dish-controls">' +
        '  <div class="dish__qty"><button data-name="' + escapeAttr(m.name) + '" data-delta="-1">−</button><span data-qty="' + escapeAttr(m.name) + '">' + inCart + '</span><button data-name="' + escapeAttr(m.name) + '" data-delta="1">+</button></div>' +
        '</div>';
      menuEl.appendChild(card);
    });
  }

  // ---- cart logic ----
  function updateQty(name, delta) {
    cart[name] = Math.max(0, (cart[name] || 0) + delta);
    if (cart[name] === 0) delete cart[name];
    renderMenu(document.querySelector('.chip.active')?.dataset.cat || '');
    renderCart();
  }

  function renderCart() {
    const names = Object.keys(cart);
    cartEmpty.style.display = names.length ? 'none' : 'block';
    cartEl.innerHTML = '';
    let total = 0;
    names.forEach(function (name) {
      const item = menu.find(function (m) { return m.name === name; });
      const qty = cart[name];
      const sub = (item ? item.price : 0) * qty;
      total += sub;
      const row = document.createElement('div');
      row.className = 'cart__item';
      row.innerHTML =
        '<div style="flex:1">' + escapeHtml(name) + ' × ' + qty + '</div>' +
        '<span class="cart__amt">' + formatRupiah(sub) + '</span>' +
        '<button data-remove="' + escapeAttr(name) + '">hapus</button>';
      cartEl.appendChild(row);
    });
    totalEl.textContent = formatRupiah(total);
    submitBtn.disabled = names.length === 0;
  }

  // ---- orders ----
  function statusText(s) {
    return { NEW: 'Diterima', PREPARING: 'Disiapkan', COMPLETED: 'Selesai', CANCELLED: 'Dibatalkan' }[s] || s;
  }
  function renderOrders(orders) {
    ordersEl.innerHTML = '';
    if (!orders.length) return;
    const head = document.createElement('div');
    head.className = 'orders__head';
    head.textContent = 'Riwayat Pesanan';
    ordersEl.appendChild(head);
    orders.forEach(function (o) {
      const items = (o.items || []).map(function (it) { return it.name + ' ×' + it.quantity; }).join(', ');
      const card = document.createElement('div');
      card.className = 'order-card';
      card.innerHTML =
        '<div class="order-card__head"><span class="order-card__no">' + escapeHtml(o.order_number) + '</span><span class="order-card__status ' + escapeHtml(o.status) + '">' + statusText(o.status) + '</span></div>' +
        '<div class="order-card__items">' + escapeHtml(items || '-') + '</div>' +
        '<div class="order-card__total">' + formatRupiah(o.total_price) + '</div>';
      ordersEl.appendChild(card);
    });
  }

  // ---- events ----
  menuEl.addEventListener('click', function (e) {
    const btn = e.target.closest('[data-delta]');
    if (btn) updateQty(btn.getAttribute('data-name'), parseInt(btn.getAttribute('data-delta'), 10));
  });
  cartEl.addEventListener('click', function (e) {
    const rm = e.target.closest('[data-remove]');
    if (rm) updateQty(rm.getAttribute('data-remove'), -9999);
  });
  catRow.addEventListener('click', function (e) {
    const chip = e.target.closest('.chip');
    if (!chip) return;
    catRow.querySelectorAll('.chip').forEach(function (c) { c.classList.remove('active'); });
    chip.classList.add('active');
    renderMenu(chip.getAttribute('data-cat') || '');
  });
  submitBtn.addEventListener('click', function () {
    const items = Object.keys(cart).map(function (name) { return { name: name, quantity: cart[name] }; });
    statusEl.textContent = 'Mengirim pesanan…';
    submitBtn.disabled = true;
    API.placeOrder(session, room, items, 'Pesanan via aplikasi kamar')
      .then(function () {
        Object.keys(cart).forEach(function (k) { delete cart[k]; });
        renderCart();
        statusEl.textContent = 'Pesanan terkirim! Mohon tunggu.';
        return API.getOrders(session);
      })
      .then(function (data) { renderOrders(data.orders || []); })
      .catch(function (err) {
        statusEl.textContent = 'Gagal mengirim: ' + (err.message || 'coba lagi');
        submitBtn.disabled = false;
      });
  });

  function escapeHtml(s) { var d = document.createElement('div'); d.textContent = s; return d.innerHTML; }
  function escapeAttr(s) { return String(s).replace(/"/g, '&quot;'); }

  // ---- init ----
  window.initDining = function () {
    document.getElementById('error').style.display = 'none';
    document.getElementById('loading').style.display = 'flex';
    API.getMenu()
      .then(function (data) { return data.items || []; })
      .then(function (items) {
        menu = items;
        document.getElementById('loading').style.display = 'none';
        renderMenu('');
        return API.getOrders(session);
      })
      .then(function (data) { renderOrders(data.orders || []); })
      .catch(function (e) {
        document.getElementById('loading').style.display = 'none';
        document.getElementById('errorText').textContent = e.message || 'Gagal memuat menu.';
        document.getElementById('error').style.display = 'flex';
      });
  };

  initDining();
})();