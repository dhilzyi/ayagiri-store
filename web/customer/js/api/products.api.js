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

export async function sendOrder(orderItems) {
  const orderID = crypto.randomUUID();
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
