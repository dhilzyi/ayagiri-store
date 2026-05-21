import { orderService, setActivePopup } from "../state.js";
import { formatterCurrency } from "./orders.ui.js";

const orderPopupItemLimit = 7;

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
  const trList = [];
  let counter = 1;
  orderService.orderList.forEach((val) => {
    const newTr = `
		<tr>
			<td>${counter}</td>
			<td>${val.product.name}</td>
			<td>${val.amount}</td>
			<td>￥${val.product.price}</td>
		</tr>
`;
    trList.push(newTr);
    counter++;
  });

  // Fill with empty row if it's still have some space
  if (trList.length < orderPopupItemLimit) {
    for (let i = trList.length; i < 7; i++) {
      trList.push(`
		<tr>
			<td>${counter}</td>
			<td></td>
			<td></td>
			<td></td>
		</tr>
`);
      counter++;
    }
  }

  // TODO: Make this fixed when the table is able to scroll
  trList.push(`
			<tr class="total-row">
				<td class="no-border"></td>
				<td class="no-border"></td>
				<td
					class="no-border"
					style="text-align: right; padding-right: 15px"
				>
					合計:
				</td>
				<td>${formatterCurrency(orderService.getTotal())}</td>
			</tr>
`);
  tbodyPopup.innerHTML = trList.join("");
}
