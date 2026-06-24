import { initListeners, startGlobalTimer } from "./ui/order_status.ui.js";
import { initSSEListen } from "./api/kitchen-api.js";
import { initDatabase } from "./handlers/database.js";
import { checkAuth } from "./auth.js";

function initMainListeners() {
  document.querySelector(".main-header").addEventListener("click", (e) => {
    const id = e.target.dataset.sectionId;
    if (!id) return;
    const targetSection = document.querySelector(
      `section[data-section-id="${id}"]`,
    );
    const activatedSection = document.querySelector("section.active");
    if (activatedSection.dataset.sectionId == id) return;

    activatedSection.classList.remove("active");
    targetSection.classList.add("active");
    console.log("changed");
  });
}

function init() {
  initMainListeners();
  initListeners();
  initDatabase();
  initSSEListen();
  startGlobalTimer();
}

const authed = await checkAuth();
if (authed) {
  document.body.style.display = "block";
  init();
} else {
  window.location.href = "/kitchen/login.html";
}
