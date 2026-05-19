import { orderService, productCache } from "../state.js";
import { renderCostTotalOrder, renderOrderItem } from "../ui/orders.ui.js";

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
    renderCostTotalOrder();
  });
}
