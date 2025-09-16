// Footer year, copy email, smooth in-page scroll
export function initFooter() {
  // Footer year
  const yearEl = document.getElementById("year");
  if (yearEl) yearEl.textContent = new Date().getFullYear();

  // Copy contact email
  const copyBtn = document.querySelector("[data-copy='#email']");
  if (copyBtn) {
    copyBtn.addEventListener("click", () => {
      const el = document.getElementById("email");
      if (!el) return;
      const text =
        el.getAttribute("href")?.replace("mailto:", "") ||
        el.textContent ||
        "";
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

      const prefersReduced = window.matchMedia(
        "(prefers-reduced-motion: reduce)"
      ).matches;

      target.scrollIntoView({
        behavior: prefersReduced ? "auto" : "smooth",
        block: "start",
      });

      history.replaceState(null, "", id);
    });
  });
}

