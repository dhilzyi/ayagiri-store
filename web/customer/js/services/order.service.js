// responsible for handling orders

export class OrderService {
  constructor() {
    this.orderList = new Map();
  }

  addToOrder(product) {
    if (this.orderList.has(product.id)) {
      this.incrementAmount(product.id);
      return true;
    }
    this.orderList.set(product.id, { product, amount: 1 });

    return true;
  }

  removeFromOrder(productId) {
    this.orderList.delete(productId);
  }

  incrementAmount(productId) {
    const item = this.orderList.get(productId);
    if (!item) {
      throw Error("item does not exist in orderList");
    }
    item.amount += 1;
  }

  decrementAmount(productId) {
    const item = this.orderList.get(productId);
    if (!item) {
      throw Error("item does not exist in orderList");
    }
    item.amount -= 1;

    if (item.amount < 1) {
      const deleted = this.orderList.delete(productId);
      return deleted;
    } else {
      return false;
    }
  }

  getTotal() {
    let total = 0;
    for (const item of this.orderList.values()) {
      total += item.product.price * item.amount;
    }
    return total;
  }

  deleteAllOrder() {
    this.orderList = new Map();
  }
}
