import { initFooter } from "./modules/footer.js";
import { initWorkAccordion } from "./modules/work.js";
import { initCarousels } from "./modules/carousel.js";
import { initThemeToggle } from "./modules/theme.js";

// Modules are deferred by default; DOM is parsed by the time this runs.
initFooter();
initWorkAccordion();
initCarousels();
initThemeToggle();

