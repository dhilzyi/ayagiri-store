import { fetchDatabase } from "../api/kitchen-api.js";
import { addRows, btnListen, renderPaginationBar } from "./database.ui.js";

class DbController {
  constructor() {
    this.table = new Map();
    this.currentPage = 1;
    this.itemsPerPage = 10;
    this.currentTable;
  }

  addDatabase(key, value) {
    this.table[key] = value;
  }

  renderPage(table, pageNum) {
    // 1. Calculate the start index (This is exactly SQL OFFSET!)
    const startIndex = (pageNum - 1) * this.itemsPerPage;

    // 2. Calculate the end index (This is exactly OFFSET + LIMIT)
    const endIndex = Math.min(
      startIndex + this.itemsPerPage,
      this.table[table].length,
    );

    // 3. Slice the data safely
    const pageData = this.table[table].slice(startIndex, endIndex);

    // 4. Render the rows (passing the starting row number for the UI, e.g., 11)
    addRows(pageData, startIndex + 1);

    // 5. Update the current page state
    this.currentPage = pageNum;
    console.log("Current Page:", this.currentPage);
    const totalPages = Math.ceil(this.table[table].length / this.itemsPerPage);
    renderPaginationBar(totalPages, this.currentPage);
  }

  renderNextRow(table) {
    const totalPages = Math.ceil(this.table[table].length / this.itemsPerPage);

    // Boundary check: Don't go past the last page
    if (this.currentPage >= totalPages) {
      return;
    }

    this.renderPage(table, this.currentPage + 1);
  }

  renderPrevRow(table) {
    // Boundary check: Don't go before page 1
    if (this.currentPage <= 1) {
      return;
    }

    this.renderPage(table, this.currentPage - 1);
  }

  changeTable(tableName) {
    if (this.currentTable == tableName) {
      return;
    }
    if (!this.table.has(tableName)) {
      throw Error("the key does not exist in the map");
    }
    this.currentTable = tableName;
  }

  getTotalPages() {
    return Math.ceil(this.table[this.currentTable].length / this.itemsPerPage);
  }
}

export const dbControl = new DbController();

export async function initDatabase() {
  const products = await fetchDatabase("products");
  console.log(products);
  dbControl.addDatabase("products", products);
  // dbControl.changeTable("products");
  dbControl.renderPage("products", 1);
  // dbControl.renderNextRow("products");
  btnListen();
}
