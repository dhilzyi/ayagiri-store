import { sendComplete } from "../api/kitchen-api.js";
import { formatSeconds } from "./helpers.js";

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

  console.log(orderData);
  const orderHeader = `
<div class="order-header">
	<h2>卓番： ${orderData.table_id}</h2>
	<span class="time" data-created-at="${orderData.created_at}">00:00</span>
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

export function startGlobalTimer() {
  setInterval(() => {
    const now = new Date();

    // Find all timer spans on the screen
    const timerSpans = document.querySelectorAll("span.time");

    timerSpans.forEach((span) => {
      const createdAtStr = span.dataset.createdAt;
      if (!createdAtStr) return;
      const localTimeStr = createdAtStr.replace("Z", "");
      const createdAt = new Date(localTimeStr);

      // Subtracting Dates in JS returns milliseconds.
      // Divide by 1000 to get seconds, and use Math.max to prevent negative numbers
      const elapsedSeconds = Math.max(0, Math.floor((now - createdAt) / 1000));

      // Update the text!
      span.textContent = formatSeconds(elapsedSeconds);
    });
  }, 1000);
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
