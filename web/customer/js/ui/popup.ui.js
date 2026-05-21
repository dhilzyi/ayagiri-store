import { orderService, setActivePopup } from "../state.js";

export function openPopupDelConfirm(bodyHTML, onConfirm) {
  document.getElementById("overlay").classList.add("active");
  document.querySelector(".popup-body").innerHTML = bodyHTML;
  document.getElementById("del-confirm").classList.add("active");

  const confirmBtn = document.getElementById("popup-confirm-btn");

  const freshBtn = confirmBtn.cloneNode(true);
  confirmBtn.replaceWith(freshBtn);
  freshBtn.addEventListener("click", () => {
    onConfirm();
    closePopupDelConfirm();
  });
  setActivePopup("del-confirm");
}

export function closePopupDelConfirm() {
  document.getElementById("overlay").classList.remove("active");
  document.getElementsByClassName("popup-body").innerHTML = ``;
  document.getElementById("del-confirm").classList.remove("active");
  setActivePopup(null);
}

export function openPopupOrderConfirm(onConfirm) {
  document.getElementById("overlay").classList.add("active");
  document.getElementById("order-confirm").classList.add("active");
  renderPopupOrderItems();

  const orderConfirmBtn = document.getElementById("order-confirm-btn");
  const cpOrderConfirmBtn = orderConfirmBtn.cloneNode(true);
  orderConfirmBtn.replaceWith(cpOrderConfirmBtn);
  cpOrderConfirmBtn.addEventListener("click", () => {
    onConfirm();
    closePopupOrderConfirm();
  });
  setActivePopup("order-confirm");
}

export function closePopupOrderConfirm() {
  document.getElementById("overlay").classList.remove("active");
  document.getElementById("order-confirm").classList.remove("active");
  setActivePopup(null);
}

export function renderPopupOrderItems() {
  const tbodyPopup = document.querySelector(".popup-order-body tbody");
  let trList = "";
  orderService.orderList.forEach((val) => {
    const newTr = `
		<tr>
			<td>${val.product.name}</td>
			<td>${val.amount}</td>
			<td>￥${val.product.price}</td>
		</tr>
`;
    trList += newTr;
  });
  tbodyPopup.innerHTML = trList;
}
