import { orderService, productCache } from "../state.js";
import {
  deleteAllOrdersItem,
  deleteOrderItem,
  renderCostTotalOrder,
  renderOrderItem,
} from "../ui/orders.ui.js";
import { closePopup, openPopup } from "../ui/popup.ui.js";

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
      renderCostTotalOrder();
    });
}

export function initPopupListeners() {
  document.getElementById("delete-all-btn").addEventListener("click", () => {
    const b = `
<p>注文をすべてキャンセルしますか？</p>
`;
    openPopup(b, () => {
      orderService.deleteAllOrder();
      deleteAllOrdersItem();
      renderCostTotalOrder();
    });
  });

  document
    .getElementById("popup-cancel-btn")
    .addEventListener("click", closePopup);
  document.getElementById("overlay").addEventListener("click", closePopup); // click outside to close
}
