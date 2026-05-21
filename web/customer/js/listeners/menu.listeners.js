import { orderService, productCache } from "../state.js";
import { renderProductsByCategory } from "../ui/menu.ui.js";
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

  document.querySelector("div.category-bar").addEventListener("click", (e) => {
    const categoryCard = e.target.closest("[data-category-id]");
    if (!categoryCard || categoryCard.classList.length > 1) return;

    document
      .querySelectorAll("[data-category-id]")
      .forEach((c) => c.classList.remove("active"));
    categoryCard.classList.add("active");
    const categoryId = Number(categoryCard.getAttribute("data-category-id"));
    renderProductsByCategory(categoryId);
  });
}
