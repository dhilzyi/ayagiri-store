import { OrderService } from "./services/order.service.js";

export const productCache = new Map();
export const orderService = new OrderService();

export function cacheProducts(products) {
  products.forEach((p) => productCache.set(p.id, p));

  console.log(productCache);
}

export function debugPrint(obj) {
  console.log(JSON.stringify(Object.fromEntries(obj), null, 2));
}
