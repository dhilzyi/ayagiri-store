// Render only module

export function renderProduct(resp) {
  const productDisplay = document.querySelector("div.product-display");
  const divEle = document.createElement("div");

  divEle.dataset.productId = 1;
  divEle.className = "product-item";
  divEle.innerHTML = `
    <h3>たこやき</h3>
    <h4>￥600<span>（税込）</span></h4>
  `;

  productDisplay.appendChild(divEle);
}
