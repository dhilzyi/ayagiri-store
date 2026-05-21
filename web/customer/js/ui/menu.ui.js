import { productCache } from "../state.js";

export function renderProduct(product) {
  const productDisplay = document.querySelector("div.product-display");
  const productItemDiv = document.createElement("div");

  productItemDiv.dataset.productId = product["id"];
  productItemDiv.className = "product-item";
  productItemDiv.innerHTML = `
    <h3>${product["name"]}</h3>
    <h4>￥${product["price"]}<span>（税込）</span></h4>
  `;

  productDisplay.appendChild(productItemDiv);
}

export function renderProductsByCategory(categoryId) {
  const productDisplay = document.querySelector("div.product-display");
  productDisplay.innerHTML = "";

  productCache.forEach((product) => {
    if (product.category_id === categoryId) {
      renderProduct(product);
    }
  });
}
