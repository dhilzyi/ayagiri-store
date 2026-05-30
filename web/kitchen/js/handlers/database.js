import { fetchDatabase } from "../api/kitchen-api.js";
import {
  addRows,
  initBtnListen,
  initPopup,
  initSelect,
  renderInformation,
  renderPaginationBar,
  updateHeader,
} from "./database.ui.js";

class DbController {
  constructor() {
    this.table = new Map();
    this.currentPage = 1;
    this.itemsPerPage;
    this.currentTable;
    this.currentHeader;
  }

  addDatabase(key, value) {
    this.table[key] = value;
  }

  addRowToTable(data) {
    console.log(data);
    this.table[this.currentTable].push(data);
    console.log(this.table[this.currentTable]);
  }

  renderPage(pageNum) {
    // Calculate the start index
    const startIndex = (pageNum - 1) * this.itemsPerPage;

    // Calculate the end index
    const endIndex = Math.min(
      startIndex + this.itemsPerPage,
      this.table[this.currentTable].length,
    );

    // Slice the data safely
    const pageData = this.table[this.currentTable].slice(startIndex, endIndex);

    // Render the rows (passing the starting row number for the UI, e.g., 11)
    addRows(pageData, startIndex + 1, this.currentTable);

    // Update the current page state
    this.currentPage = pageNum;
    console.log("Current Page:", this.currentPage);
    const totalPages = Math.ceil(
      this.table[this.currentTable].length / this.itemsPerPage,
    );
    renderPaginationBar(totalPages, this.currentPage);
    renderInformation(
      this.table[this.currentTable].length,
      startIndex + 1,
      endIndex,
    );
  }

  renderNextRow() {
    const totalPages = Math.ceil(
      this.table[this.currentTable].length / this.itemsPerPage,
    );

    // Boundary check: Don't go past the last page
    if (this.currentPage >= totalPages) {
      return;
    }

    this.renderPage(this.currentPage + 1);
  }

  renderPrevRow() {
    // Boundary check: Don't go before page 1
    if (this.currentPage <= 1) {
      return;
    }

    this.renderPage(this.currentPage - 1);
  }

  // Change Table Database
  async changeTable(tableName) {
    if (this.currentTable == tableName) {
      return;
    }
    if (!Object.hasOwn(this.table, tableName)) {
      try {
        const data = await fetchDatabase(tableName);
        this.addDatabase(tableName, data);
      } catch (err) {
        throw err;
      }
    }

    this.currentTable = tableName;
    this.currentHeader = tableName;
    updateHeader(tableName);
    this.renderPage(1);
  }

  getTotalPages() {
    return Math.ceil(this.table[this.currentTable].length / this.itemsPerPage);
  }

  // async settting items per page
  async setItemsPerPage(value) {
    this.itemsPerPage = value;
  }
}

export const dbControl = new DbController();

export async function initDatabase() {
  // const products = await fetchDatabase("products");
  // dbControl.addDatabase("products", products);
  // dbControl.changeTable("products");
  // await dbControl.changeTable("products");
  // dbControl.renderNextRow("products");
  await initSelect();
  initBtnListen();
  initPopup();
}
