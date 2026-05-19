import { orderService, productCache } from "../state.js";
import {
  deleteOrderItem,
  renderCostTotal,
  renderOrderItem,
} from "../ui/orders.ui.js";

export function initMenuListeners() {
  document.querySelector(".product-display").addEventListener("click", (e) => {
    const card = e.target.closest("[data-product-id]");
    if (!card) return;

    const productID = Number(card.getAttribute("data-product-id"));
    const product = productCache.get(productID);
    if (!product) {
      throw Error("product with following ID is not found in map");
    }
    const added = orderService.addToOrder(product);
    if (added) {
      renderOrderItem(product);
    }
    renderCostTotal();
  });
}

export function initOrderListeners() {
  document
    .querySelector("div.item-lists table")
    .addEventListener("click", (e) => {
      const tr = e.target.closest("tr[data-product-id]");
      if (!tr) return;

      const productID = Number(tr.dataset.productId);
      const product = productCache.get(productID);
      if (!product)
        throw Error("product with following ID is not found in map");

      if (e.target.closest(".btn-increment")) {
        orderService.incrementAmount(productID);
        renderOrderItem(product);
      } else if (e.target.closest(".btn-decrement")) {
        const deleted = orderService.decrementAmount(productID);
        if (deleted) {
          deleteOrderItem(productID);
        } else {
          renderOrderItem(product);
        }
      }
      renderCostTotal();
    });
}
