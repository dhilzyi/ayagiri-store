import { fetchDatabase, sendDeleteRows } from "../api/kitchen-api.js";
import {
  addRows,
  renderInformation,
  renderPaginationBar,
  updateHeader,
} from "./database.ui.js";

import {
  getSelectedIDs,
  initBtnListen,
  initPopup,
  initTableDbSelect,
} from "./database.events.js";

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
    this.table[this.currentTable].push(data);
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

    if (this.currentPage >= totalPages) {
      return;
    }

    this.renderPage(this.currentPage + 1);
  }

  renderPrevRow() {
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

  async deleteRows() {
    const idsToDelete = getSelectedIDs();

    if (idsToDelete.length === 0) return;

    try {
      await sendDeleteRows(idsToDelete, this.currentTable);
    } catch (err) {
      return;
    }

    const activeTableData = this.table[this.currentTable];
    const updatedData = activeTableData.filter((item) => {
      return !idsToDelete.includes(item.id);
    });

    this.table[this.currentTable] = updatedData;

    if (this.currentPage > this.getTotalPages() && this.currentPage > 1) {
      this.currentPage = this.getTotalPages();
    }

    this.renderPage(this.currentPage);
  }

  getDataByID(id) {
    const index = binarySearch(this.table[this.currentTable], id);
    if (index === -1) throw Error("no id is found in the table");
    return this.table[this.currentTable][index];
  }
  async updateDataByID(id, data) {
    const index = binarySearch(this.table[this.currentTable], id);
    if (index === -1) throw Error("no id is found in the table");
    const current = this.table[this.currentTable][index];
    current.category_id = data.category_id;
    current.discount = data.discount;
    current.name = data.name;
    current.price = data.price;
  }
}

// It is guarantee the table is always in order because i use the query ASC
function binarySearch(arr, val) {
  let start = 0;
  let end = arr.length - 1;

  while (start <= end) {
    let mid = Math.floor((start + end) / 2);

    if (arr[mid].id === val) {
      return mid;
    }

    if (val < arr[mid].id) {
      end = mid - 1;
    } else {
      start = mid + 1;
    }
  }
  return -1;
}

export const dbControl = new DbController();

export async function initDatabase() {
  await initTableDbSelect();
  initBtnListen();
  initPopup();
}
