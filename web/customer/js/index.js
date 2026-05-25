import { loadProducts } from "./api/customers.api.js";
import { initListeners } from "./listeners/index.js";
import { cacheProducts } from "./state.js";
import { renderProductsByCategory } from "./ui/menu.ui.js";

const products = await loadProducts();

cacheProducts(products);

// initial render
renderProductsByCategory(2);

initListeners();
