import { initFooter } from "./modules/footer.js";
import { initProjectsAccordion } from "./modules/projects.js";
import { initCarousels } from "./modules/carousel.js";
import { initThemeToggle } from "./modules/theme.js";

// Modules are deferred by default; DOM is parsed by the time this runs.
initFooter();
initProjectsAccordion();
initCarousels();
initThemeToggle();

