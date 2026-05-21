import { initPopupListeners, initOrderListeners } from "./order.listeners.js";
import { initMenuListeners } from "./menu.listeners.js";

export function initListeners() {
  initMenuListeners();
  initOrderListeners();
  initPopupListeners();
}
