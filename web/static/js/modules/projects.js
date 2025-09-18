// web/static/js/modules/projects.js

/**
 * Handles the interactive accordion behavior for the "Projects" section.
 * - Ensures only one project card can be open at a time.
 * - On small screens, scrolls the opened card into view for better visibility.
 * - Toggles a `has-open` class on the container for advanced CSS styling.
 */
export function initProjectsAccordion() {
  // Find the project stack using the correct ID: #projects
  const stack = document.querySelector("#projects .project-stack");
  if (!stack) {
    // Exit gracefully if the projects section isn't on the page.
    return;
  }

  /**
   * Toggles a 'has-open' class on the project stack container.
   * This class is used by CSS to switch between grid layouts.
   */
  const updateHasOpenClass = () => {
    const isAnyCardOpen = !!stack.querySelector(".project-card[open]");
    stack.classList.toggle("has-open", isAnyCardOpen);
  };

  // Listen for the 'toggle' event, which fires whenever a <details> element is opened or closed.
  stack.addEventListener(
    "toggle",
    (event) => {
      const activeCard = event.target;
      // Ensure the event came from a project card.
      if (!activeCard.matches(".project-card")) {
        return;
      }

      // If a card was just opened, close all other open cards.
      if (activeCard.open) {
        for (const card of stack.querySelectorAll("details.project-card[open]")) {
          if (card !== activeCard) {
            card.open = false;
          }
        }

        // On smaller screens, scroll the active card to the top for a better user experience.
        if (window.innerWidth < 1024) {
          activeCard.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      }

      // After any toggle, update the parent container's class.
      updateHasOpenClass();
    },
    true // Use event capturing to ensure this runs before other potential listeners.
  );

  // Run once on initialization to set the correct initial state.
  updateHasOpenClass();
}

