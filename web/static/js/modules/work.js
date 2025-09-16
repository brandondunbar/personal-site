// "Selected Work" accordion (only one open, scroll into view on small screens)
export function initWorkAccordion() {
  const stack = document.querySelector("#work .project-stack");
  if (!stack) return;

  const updateHasOpenClass = () => {
    const anyOpen = !!stack.querySelector(".project-card[open]");
    stack.classList.toggle("has-open", anyOpen);
  };

  stack.addEventListener(
    "toggle",
    (e) => {
      const el = e.target;
      if (el.tagName !== "DETAILS") return;

      if (el.open) {
        for (const d of stack.querySelectorAll("details.project-card[open]")) {
          if (d !== el) d.open = false;
        }
        if (window.innerWidth < 1024) {
          el.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      }
      updateHasOpenClass();
    },
    true
  );

  updateHasOpenClass();
}

