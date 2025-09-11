// Fill current year in footer
document.addEventListener("DOMContentLoaded", () => {
  const yearEl = document.getElementById("year");
  if (yearEl) yearEl.textContent = new Date().getFullYear();

  // Copy-to-clipboard for contact email
  const copyBtn = document.querySelector("[data-copy='#email']");
  if (copyBtn) {
    copyBtn.addEventListener("click", () => {
      const a = document.getElementById("email");
      if (!a) return;
      const text = a.getAttribute("href")?.replace("mailto:", "") || a.textContent;
      navigator.clipboard.writeText(text).then(() => {
        const original = copyBtn.textContent;
        copyBtn.textContent = "Copied ✓";
        setTimeout(() => (copyBtn.textContent = original), 1200);
      });
    });
  }

  // Smooth scroll for in-page anchors (reduced-motion friendly)
  document.querySelectorAll('a[href^="#"]').forEach(a => {
    a.addEventListener("click", e => {
      const id = a.getAttribute("href");
      if (!id || id === "#") return;
      const el = document.querySelector(id);
      if (!el) return;
      e.preventDefault();
      const prefersReduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      el.scrollIntoView({ behavior: prefersReduced ? "auto" : "smooth", block: "start" });
      history.replaceState(null, "", id);
    });
  });
});

// ---- Theme toggle ----
(function(){
  const STORAGE_KEY = 'theme'; // 'light' | 'dark' | null
  const html = document.documentElement;
  const btn = document.getElementById('themeToggle');
  if(!btn) return;

  const iconSun  = btn.querySelector('.icon-sun');
  const iconMoon = btn.querySelector('.icon-moon');

  function applyTheme(t){
    html.classList.remove('theme-light','theme-dark');
    if (t === 'light') html.classList.add('theme-light');
    else if (t === 'dark') html.classList.add('theme-dark');
    // update button visuals
    const dark = getComputedStyle(document.documentElement).getPropertyValue('--bg').trim().startsWith('#0') || html.classList.contains('theme-dark');
    btn.setAttribute('aria-pressed', dark ? 'true' : 'false');
    if(dark){ iconSun.style.display='none'; iconMoon.style.display='inline'; }
    else    { iconSun.style.display='inline'; iconMoon.style.display='none'; }
  }

  // Initial theme: stored > system
  const stored = localStorage.getItem(STORAGE_KEY);
  if(stored){ applyTheme(stored); }
  else {
    const prefersLight = window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches;
    applyTheme(prefersLight ? 'light' : 'dark');
  }

  btn.addEventListener('click', ()=>{
    const next = html.classList.contains('theme-light') ? 'dark' : 'light';
    localStorage.setItem(STORAGE_KEY, next);
    applyTheme(next);
  });
})();

// ---- Subtle parallax on .hero-graphic ----
(function(){
  const el = document.querySelector('.hero-graphic');
  if(!el) return;
  const speed = 0.12; // lower = subtler

  function update(){
    const rect = el.getBoundingClientRect();
    const vh = window.innerHeight || document.documentElement.clientHeight;
    // Only animate when on screen
    if(rect.bottom > 0 && rect.top < vh){
      const centerOffset = (rect.top + rect.height/2) - vh/2;
      const y = Math.round(-centerOffset * speed);
      el.style.setProperty('--parallax-y', y + 'px');
    }
  }

  let ticking = false;
  function onScroll(){
    if(!ticking){
      window.requestAnimationFrame(()=>{ update(); ticking=false; });
      ticking = true;
    }
  }
  update();
  window.addEventListener('scroll', onScroll, {passive:true});
  window.addEventListener('resize', update);
})();

(() => {
  const STORAGE_KEY = "theme"; // "light" | "dark"
  const root = document.documentElement;
  const btn = document.getElementById("theme-toggle");

  // Decide initial theme: saved -> system -> default "dark"
  const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const saved = localStorage.getItem(STORAGE_KEY);
  const initial = (saved === "light" || saved === "dark") ? saved : (systemPrefersDark ? "dark" : "light");
  applyTheme(initial, { noTransition: true });

  // Mark ready to prevent icon flash
  root.classList.add("theme-ready");

  // Toggle on click
  btn?.addEventListener("click", () => {
    const next = root.dataset.theme === "dark" ? "light" : "dark";
    applyTheme(next);
  });

  function applyTheme(mode, opts = {}) {
    if (!opts.noTransition) enableTransitionOnce();
    root.setAttribute("data-theme", mode);
    // Persist explicit user choices
    if (!opts.systemDriven) localStorage.setItem(STORAGE_KEY, mode);
    // Keep the button label helpful for screen readers
    const label = mode === "dark" ? "Switch to light theme" : "Switch to dark theme";
    btn?.setAttribute("aria-label", label);
    btn?.setAttribute("title", label);
  }

  function enableTransitionOnce() {
    root.classList.add("theme-transition");
    window.setTimeout(() => root.classList.remove("theme-transition"), 170);
  }
})();

