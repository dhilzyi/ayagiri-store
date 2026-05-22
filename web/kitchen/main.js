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

initMainListeners();
