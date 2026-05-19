export function openPopup(bodyHTML, onConfirm) {
  document.getElementById("overlay").classList.add("active");
  document.querySelector(".popup-body").innerHTML = bodyHTML;
  document.getElementById("popup-confirm").classList.add("active");

  const confirmBtn = document.getElementById("popup-confirm-btn");

  const freshBtn = confirmBtn.cloneNode(true);
  confirmBtn.replaceWith(freshBtn);
  freshBtn.addEventListener("click", () => {
    onConfirm();
    closePopup();
  });
}

export function closePopup() {
  document.getElementById("overlay").classList.remove("active");
  document.getElementsByClassName("popup-body").innerHTML = ``;
  document.getElementById("popup-confirm").classList.remove("active");
}
