import { loadProductsByCategory } from "./api/products.api.js";
import { initListeners } from "./listeners/index.js";
import { cacheProducts } from "./state.js";
import { renderProduct } from "./ui/products.ui.js";

const products = await loadProductsByCategory("1");
products.forEach((product) => renderProduct(product));
cacheProducts(products);

initListeners();
