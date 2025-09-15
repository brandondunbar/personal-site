// ------------------------------
// Footer year + copy to clipboard
// ------------------------------
document.addEventListener("DOMContentLoaded", () => {
  // Footer year
  const yearEl = document.getElementById("year");
  if (yearEl) yearEl.textContent = new Date().getFullYear();

  // Copy contact email
  const copyBtn = document.querySelector("[data-copy='#email']");
  if (copyBtn) {
    copyBtn.addEventListener("click", () => {
      const el = document.getElementById("email");
      if (!el) return;
      const text = el.getAttribute("href")?.replace("mailto:", "") || el.textContent || "";
      navigator.clipboard.writeText(text).then(() => {
        const original = copyBtn.textContent;
        copyBtn.textContent = "Copied ✓";
        setTimeout(() => (copyBtn.textContent = original), 1200);
      });
    });
  }

  // Smooth scroll for in-page anchors (reduced-motion friendly)
  document.querySelectorAll('a[href^="#"]').forEach((a) => {
    a.addEventListener("click", (e) => {
      const id = a.getAttribute("href");
      if (!id || id === "#") return;
      const target = document.querySelector(id);
      if (!target) return;
      e.preventDefault();
      const prefersReduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      target.scrollIntoView({ behavior: prefersReduced ? "auto" : "smooth", block: "start" });
      history.replaceState(null, "", id);
    });
  });
});

// ------------------------------
// Theme toggle (single, robust)
// ------------------------------
(() => {
  const STORAGE_KEY = "theme"; // "light" | "dark"
  const root = document.documentElement;

  // Support either #theme-toggle (new) or #themeToggle (old)
  const btn = document.getElementById("theme-toggle") || document.getElementById("themeToggle");
  if (!btn) return;

  // Decide initial theme: saved -> system -> default "dark"
  const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const saved = localStorage.getItem(STORAGE_KEY);
  const initial = saved === "light" || saved === "dark" ? saved : systemPrefersDark ? "dark" : "light";
  applyTheme(initial, { noTransition: true, systemDriven: !saved });

  // Mark ready to prevent icon flash (CSS uses .theme-ready)
  root.classList.add("theme-ready");

  // Toggle on click
  btn.addEventListener("click", () => {
    const next = root.dataset.theme === "dark" ? "light" : "dark";
    applyTheme(next);
  });

  function applyTheme(mode, opts = {}) {
    if (!opts.noTransition) enableTransitionOnce();
    root.setAttribute("data-theme", mode);
    if (!opts.systemDriven) localStorage.setItem(STORAGE_KEY, mode);

    const label = mode === "dark" ? "Switch to light theme" : "Switch to dark theme";
    btn.setAttribute("aria-label", label);
    btn.setAttribute("title", label);
    btn.setAttribute("aria-pressed", mode === "dark" ? "true" : "false");
  }

  function enableTransitionOnce() {
    root.classList.add("theme-transition");
    window.setTimeout(() => root.classList.remove("theme-transition"), 170);
  }
})();

// --------------------------------------
// Selected Work: one-open + scroll/fallback
// --------------------------------------
(() => {
  const stack = document.querySelector("#work .project-stack");
  if (!stack) return;

  // Helper to update fallback class for environments without :has()
  const updateHasOpenClass = () => {
    const anyOpen = !!stack.querySelector(".project-card[open]");
    stack.classList.toggle("has-open", anyOpen);
  };

  stack.addEventListener(
    "toggle",
    (e) => {
      const el = e.target;
      if (el.tagName !== "DETAILS") return;

      // If this card just opened, close the others (accordion behavior)
      if (el.open) {
        for (const d of stack.querySelectorAll("details.project-card[open]")) {
          if (d !== el) d.open = false;
        }
        // On narrow screens, ensure the opened card is visible
        if (window.innerWidth < 1024) {
          el.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      }

      // Maintain a CSS fallback state for browsers without :has()
      updateHasOpenClass();
    },
    true
  );

  // Initialize fallback class on load
  updateHasOpenClass();
})();


