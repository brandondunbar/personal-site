// web/static/js/modules/carousel.js

function whenImagesReady(container, cb) {
  const imgs = Array.from(container.querySelectorAll('img'));
  if (imgs.length === 0) { cb(); return; }
  let remaining = imgs.length;
  const done = () => { if (--remaining <= 0) cb(); };
  imgs.forEach(img => {
    if (img.complete) { requestAnimationFrame(done); }
    else {
      img.addEventListener('load', done,  { once: true });
      img.addEventListener('error', done, { once: true });
    }
  });
}

function pickScroller(root, track) {
  // Prefer whichever actually overflows horizontally
  if (track.scrollWidth > track.clientWidth + 1) return track;
  if (root.scrollWidth  > root.clientWidth  + 1) return root;
  return (root.scrollWidth > track.scrollWidth) ? root : track;
}

function attachCarousel(targetId) {
  const rootEl  = document.getElementById(targetId);           // .carousel
  if (!rootEl) return;
  const trackEl = rootEl.querySelector('.carousel__track');     // <ul>
  if (!trackEl) return;

  const prevBtn = document.querySelector(`[data-carousel-prev="${targetId}"]`);
  const nextBtn = document.querySelector(`[data-carousel-next="${targetId}"]`);

  let scroller = pickScroller(rootEl, trackEl);

  const setHidden = (btn, hidden) => {
    if (!btn) return;
    btn.classList.toggle('is-disabled', hidden);
    btn.setAttribute('aria-hidden', hidden ? 'true' : 'false');
    btn.tabIndex = hidden ? -1 : 0;
  };

  const measure = () => {
    scroller = pickScroller(rootEl, trackEl);
    const maxScroll = Math.max(0, scroller.scrollWidth - scroller.clientWidth - 1);
    const atStart   = scroller.scrollLeft <= 0;
    const atEnd     = scroller.scrollLeft >= maxScroll;
    const noOverflow = maxScroll <= 0;
    setHidden(prevBtn, noOverflow || atStart);
    setHidden(nextBtn, noOverflow || atEnd);
  };

  const stepSize = () => {
    const card = trackEl.querySelector('.carousel__slide');
    const cs   = getComputedStyle(trackEl);
    const gap  = parseFloat(cs.columnGap || cs.gap || 16);
    const width = card ? card.clientWidth : 144;
    return width + gap;
  };

  const scrollByCards = (n) => {
    scroller.scrollBy({ left: stepSize() * n, behavior: 'smooth' });
  };

  prevBtn?.addEventListener('click', () => scrollByCards(-3));
  nextBtn?.addEventListener('click', () => scrollByCards(+3));

  const onScroll = () => measure();
  const rebindScrollListener = () => {
    scroller.removeEventListener?.('scroll', onScroll);
    scroller = pickScroller(rootEl, trackEl);
    scroller.addEventListener('scroll', onScroll, { passive: true });
  };

  rebindScrollListener();
  window.addEventListener('resize', () => { rebindScrollListener(); measure(); });

  // Keyboard when carousel has focus
  rootEl.addEventListener('keydown', (e) => {
    if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
      e.preventDefault();
      scrollByCards(e.key === 'ArrowRight' ? +1 : -1);
    }
  });

  const ro = new ResizeObserver(() => { rebindScrollListener(); measure(); });
  ro.observe(rootEl);
  ro.observe(trackEl);

  const mo = new MutationObserver(() => { rebindScrollListener(); measure(); });
  mo.observe(trackEl, { childList: true, subtree: true });

  whenImagesReady(trackEl, () => { rebindScrollListener(); measure(); });
  requestAnimationFrame(() => { rebindScrollListener(); measure(); });
  window.addEventListener('load', () => { rebindScrollListener(); measure(); });
}

export function initCarousels() {
  document.querySelectorAll('.js-carousel').forEach(c => attachCarousel(c.id));
}

