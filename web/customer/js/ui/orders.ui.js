import { orderService } from "../state.js";

export function renderOrderItem(product) {
  const tableOrders = document.querySelector(".item-lists table");
  let trItem = tableOrders.querySelector(`tr[data-product-id='${product.id}']`);
  if (trItem) {
    const item = orderService.orderList.get(product.id);
    trItem.querySelector(".product-amount").textContent = item.amount;
  } else {
    trItem = document.createElement("tr");
    trItem.dataset.productId = product.id;
    trItem.innerHTML = `
	<td class="product-name">${product.name}</td>
	<td class="product-amount">1</td>
	<td>
		<button class="btn-decrement">-</button>
		<button class="btn-increment">+</button>
	</td>
`;
    tableOrders.appendChild(trItem);
  }
}

export function deleteOrderItem(productId) {
  const trItem = document.querySelector(`tr[data-product-id='${productId}']`);
  trItem.remove();
}

export function renderCostTotal() {
  const h2Price = document.querySelector("h2.sum-price");
  h2Price.textContent = `￥${orderService.getTotal()}`;
}
