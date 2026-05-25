import { initSSEOrders, sendOrder } from "../api/customers.api.js";
import { activePopup, orderService, productCache } from "../state.js";
import {
  deleteAllOrdersItem,
  deleteOrderItem,
  renderCostTotalOrder,
  renderOrderItem,
} from "../ui/orders.ui.js";
import {
  closePopupDelConfirm,
  closePopupOrderConfirm,
  openPopupDelConfirm,
  openPopupOrderConfirm,
} from "../ui/popup.ui.js";

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
  // Overlay cancel button listener
  document.getElementById("overlay").addEventListener("click", () => {
    if (activePopup === "del-confirm") {
      closePopupDelConfirm();
    }
    if (activePopup === "order-confirm") {
      closePopupOrderConfirm();
    }
  });

  // Delete orders confirm popup, accept/cancel button listener
  document.getElementById("delete-all-btn").addEventListener("click", () => {
    const b = `
<p>注文をすべてキャンセルしますか？</p>
`;
    openPopupDelConfirm(b, () => {
      orderService.deleteAllOrder();
      deleteAllOrdersItem();
    });
  });
  document
    .getElementById("popup-cancel-btn")
    .addEventListener("click", closePopupDelConfirm);

  // Order confirm popup. accept/cancel button listener
  document.getElementById("order-btn").addEventListener("click", () => {
    openPopupOrderConfirm(async () => {
      console.log("Sending order...");

      const orderItems = orderService.getArrayOrderItems();
      const orderID = crypto.randomUUID();
      try {
        await sendOrder(orderItems, orderID);
      } catch (err) {
        console.log(err);
        return;
      }
      initSSEOrders(orderID);
      orderService.deleteAllOrder();
      deleteAllOrdersItem();
    });
  });
  document
    .getElementById("order-confirm-cancel-btn")
    .addEventListener("click", closePopupOrderConfirm);
}
