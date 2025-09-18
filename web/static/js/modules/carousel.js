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

function attachCarousel(targetId) {
  const rootEl  = document.getElementById(targetId);           // .carousel (the ONLY scroller)
  if (!rootEl) return;
  const trackEl = rootEl.querySelector('.carousel__track');     // <ul> (content row)
  if (!trackEl) return;

  const prevBtn = document.querySelector(`[data-carousel-prev="${targetId}"]`);
  const nextBtn = document.querySelector(`[data-carousel-next="${targetId}"]`);

  const setHidden = (btn, hidden) => {
    if (!btn) return;
    btn.classList.toggle('is-disabled', hidden);
    btn.setAttribute('aria-hidden', hidden ? 'true' : 'false');
    btn.tabIndex = hidden ? -1 : 0;
  };

  // Small tolerance to handle subpixel scroll values and rounding
  const EPS = 1.5;

  const measure = () => {
    const maxScroll = Math.max(0, rootEl.scrollWidth - rootEl.clientWidth - 1);
    const atStart   = rootEl.scrollLeft <= EPS;
    const atEnd     = rootEl.scrollLeft >= (maxScroll - EPS);
    const noOverflow = maxScroll <= EPS;

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
    rootEl.scrollBy({ left: stepSize() * n, behavior: 'smooth' });
  };

  // Click handlers
  prevBtn?.addEventListener('click', () => scrollByCards(-3));
  nextBtn?.addEventListener('click', () => scrollByCards(+3));

  // Keep arrow state in sync
  const onScroll = () => measure();
  rootEl.addEventListener('scroll', onScroll, { passive: true });
  window.addEventListener('resize', measure);

  // Keyboard support when the carousel has focus
  rootEl.addEventListener('keydown', (e) => {
    if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
      e.preventDefault();
      scrollByCards(e.key === 'ArrowRight' ? +1 : -1);
    }
  });

  // Make vertical wheel/trackpad gestures scroll horizontally (while scrollbar is hidden)
  rootEl.addEventListener('wheel', (e) => {
    // Ignore pinch-zoom (ctrlKey) and gestures already mostly horizontal
    if (e.ctrlKey) return;
    const absX = Math.abs(e.deltaX);
    const absY = Math.abs(e.deltaY);
    if (absY > absX && absY > 0) {
      // Translate vertical delta to horizontal scroll
      rootEl.scrollBy({ left: e.deltaY, behavior: 'auto' });
      e.preventDefault(); // allow smooth native feel; keep passive: false
    }
  }, { passive: false });

  // Recompute on size/content changes
  const ro = new ResizeObserver(() => measure());
  ro.observe(rootEl);
  ro.observe(trackEl);

  const mo = new MutationObserver(() => measure());
  mo.observe(trackEl, { childList: true, subtree: true });

  // Wait for images, then measure; also measure on next frame and on window load
  whenImagesReady(trackEl, measure);
  requestAnimationFrame(measure);
  window.addEventListener('load', measure);

  // Optional cleanup handle if needed later
  return () => {
    prevBtn?.removeEventListener('click', scrollByCards);
    nextBtn?.removeEventListener('click', scrollByCards);
    rootEl.removeEventListener('scroll', onScroll);
    window.removeEventListener('resize', measure);
    window.removeEventListener('load', measure);
    ro.disconnect();
    mo.disconnect();
  };
}

export function initCarousels() {
  document.querySelectorAll('.js-carousel').forEach(c => attachCarousel(c.id));
}

