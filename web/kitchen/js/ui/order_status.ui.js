import { sendComplete } from "../api/kitchen-api.js";

export function renderOrderTotals() {
  const orderList = document.querySelectorAll(".order-card");
  const orderNumber = document
    .querySelector(".summary-box")
    .firstElementChild.querySelector("span");
  orderNumber.textContent = orderList.length;

  console.log(orderList.length);
  const orderItems = document.querySelectorAll("div.item");
  console.log(orderItems.length);
}

export function addOrderToQueue(orderData) {
  const orderList = document.querySelector("div.order-list");
  const itemList = [];
  orderData.items.forEach((obj) => {
    itemList.push(`
	<div class="item">
		<span>${obj.products.name}</span><span>x ${obj.quantity}</span>
	</div>
`);
  });

  const orderHeader = `
<div class="order-header">
	<h2>卓番： ${orderData.table_id}</h2>
	<span class="time">00:00</span>
</div>
`;

  const orderContent = `
<div class="order-content">
	<div class="order-items-box">
		${itemList.join("")}
	</div>
	<button class="complete-btn">
		完成
	</button>
</div>
`;
  const orderCard = document.createElement("div");
  orderCard.classList.add("order-card");
  orderCard.dataset.orderId = orderData.order_id;
  orderCard.innerHTML = orderHeader + orderContent;

  orderList.appendChild(orderCard);
}

export function removeOrderFromQueue(orderID) {
  const orderCard = document.querySelector(
    `.order-card[data-order-id='${orderID}']`,
  );
  if (!orderCard) return;
  orderCard.remove();
  renderOrderTotals();
}

export function initListeners() {
  document.querySelector(".order-list").addEventListener("click", (e) => {
    const btn = e.target.closest(".complete-btn");
    if (!btn) return;
    const orderID = e.target.closest("[data-order-id]").dataset.orderId;
    try {
      sendComplete(orderID);
      removeOrderFromQueue(orderID);
    } catch (err) {
      console.log(err);
      return;
    }
  });
}
