import { loadProducts, loadProductsByCategory } from "./api/products.api.js";
import { initListeners } from "./listeners/index.js";
import { cacheProducts } from "./state.js";
import { renderProduct, renderProductsByCategory } from "./ui/menu.ui.js";

const products = await loadProducts();

cacheProducts(products);

// initial render
renderProductsByCategory(2);

initListeners();
