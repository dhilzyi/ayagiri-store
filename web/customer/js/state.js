import { OrderService } from "./services/order.service.js";

export const productCache = new Map();
export const orderService = new OrderService();

export function cacheProducts(products) {
  products.forEach((p) => productCache.set(p.id, p));
}

export function debugPrint(obj) {
  console.log(JSON.stringify(Object.fromEntries(obj), null, 2));
}

// Popup overlay state
export let activePopup = null;
export function setActivePopup(id) {
  activePopup = id;
}
