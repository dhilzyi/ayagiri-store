// Render only module

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
