// Fetch to the api only module

export async function loadProductsByCategory(categoryID) {
  const res = await fetch(`/api/products?category_id=${categoryID}`);
  if (!res.ok) {
    throw Error(`response bad status: ${res.status}`);
  }
  const products = res.json();

  return products;
}

export async function loadProducts() {
  const res = await fetch(`/api/product`);
  if (!res.ok) {
    throw Error(`response bad status: ${res.status}`);
  }
  const products = res.json();

  return products;
}

export async function sendOrder(orderItems, orderID) {
  const payload = {
    table_id: 2,
    items: orderItems,
  };
  console.log(orderID);
  console.log(JSON.stringify(payload));
  const res = await fetch(`/api/orders?order_id=${orderID}`, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: {
      "Content-type": "application/json",
    },
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(
      errorData.message || `Server returned status ${res.status}`,
    );
  }
  console.log(res.status);
}

// TODO: Optimize to only have one SSE connection per customer client not per orders
export async function initSSEOrders(orderID) {
  const eventSource = new EventSource(`/api/orders/stream?order_id=${orderID}`);
  eventSource.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log(msg);

    eventSource.close();
    console.log("EventSource closed successfully.");
  };

  eventSource.onerror = (err) => {
    console.log(err);
  };
}
