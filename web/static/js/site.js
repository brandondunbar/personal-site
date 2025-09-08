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
