// Theme toggle: saves to localStorage, syncs ARIA
const KEY = "theme";

export function initThemeToggle() {
  const root = document.documentElement;
  const btn =
    document.getElementById("theme-toggle") ||
    document.getElementById("themeToggle");
  if (!btn) return;

  const setTheme = (t, persist = true) => {
    root.setAttribute("data-theme", t);
    if (persist) {
      try {
        localStorage.setItem(KEY, t);
      } catch { /* ignore */ }
    }
    const label = t === "dark" ? "Switch to light theme" : "Switch to dark theme";
    btn.setAttribute("aria-label", label);
    btn.setAttribute("title", label);
    btn.setAttribute("aria-pressed", String(t === "dark")); // <-- fixed
  };

  // initialize: saved -> system -> light
  let saved = null;
  try { saved = localStorage.getItem(KEY); } catch { /* ignore */ }

  if (saved === "light" || saved === "dark") {
    setTheme(saved, false);
  } else {
    const prefersDark =
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches;
    setTheme(prefersDark ? "dark" : "light", false);
  }

  btn.addEventListener("click", () => {
    const cur = root.getAttribute("data-theme") === "dark" ? "dark" : "light";
    setTheme(cur === "dark" ? "light" : "dark");
  });
}

