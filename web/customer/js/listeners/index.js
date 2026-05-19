import { initPopupListeners, initOrderListeners } from "./order.listeners.js";
import { initMenuListeners } from "./products.listeners.js";

export function initListeners() {
  initMenuListeners();
  initOrderListeners();
  initPopupListeners();
}
