// Fetch to the api only module

export async function loadProductsByCategory(categoryID) {
  const res = await fetch(`/api/product?category_id=${categoryID}`);
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
