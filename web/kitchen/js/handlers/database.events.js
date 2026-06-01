import { dbControl } from "./database.js";
import { sendNewRows, sendUpdateRows } from "../api/kitchen-api.js";
import { fillPopup, renderResults } from "./database.ui.js";

export async function initTableDbSelect() {
  const databaseSelect = document.getElementById("database-select");
  const perPage = document.getElementById("result-per-page-opt");
  databaseSelect.addEventListener("change", () => {
    dbControl.changeTable(databaseSelect.value);
  });
  perPage.addEventListener("change", async () => {
    await dbControl.setItemsPerPage(Number(perPage.value));
    dbControl.renderPage(dbControl.currentPage);
  });
  await dbControl.setItemsPerPage(Number(perPage.value));
  await dbControl.changeTable(databaseSelect.value);
  // DELETE BEFORE COMMIT
  dbControl.renderPage(dbControl.getTotalPages());
}

export function initPopup() {
  const modal = document.getElementById("create-row");
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.close();
    }
  });

  const form = document.getElementById("popup-form");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const table = form.dataset.table;
    const action = form.dataset.action;
    const formData = new FormData(form);

    let data = Object.fromEntries(formData);
    data = parseDataForm(data, table);
    switch (action) {
      case "add": {
        try {
          const response = await sendNewRows(JSON.stringify(data), table);

          renderResults({ status: response.status }, "success");
          dbControl.addRowToTable(await response.json());
        } catch (err) {
          const errData = (await err.body) || {};
          console.log(err);
          errData.status = err.status || 500;
          errData.message = err.error;

          renderResults(errData, "error");
          return;
        }
        dbControl.renderPage(dbControl.getTotalPages());
        form.reset();
        break;
      }
      case "update": {
        const productId = Number(form.dataset.id);
        const response = await sendUpdateRows(
          productId,
          JSON.stringify(data),
          table,
        );
        await dbControl.updateDataByID(response.body.id, response.body);
        dbControl.renderPage(dbControl.currentPage);
        modal.close();
        form.reset();
        break;
      }
    }
  });

  // DEBUG POP UP
  // modal.showModal();

  document
    .querySelector(".action .right")
    .addEventListener("click", async (e) => {
      const cardId = e.target.id;
      if (!cardId) return;
      switch (cardId) {
        case "new-order-btn": {
          modal.showModal();
          break;
        }
        case "delete-btn": {
          dbControl.deleteRows();
          break;
        }
        case "edit-btn": {
          const toEditId = getSelectedIDs();
          const rowData = dbControl.getDataByID(toEditId[0]);
          await fillPopup(rowData);
          modal.showModal();

          break;
        }
      }
    });
  document.getElementById("popup-cancel-btn").addEventListener("click", () => {
    modal.close();
  });
}

export function initBtnListen() {
  document
    .querySelector(".database-table thead")
    .addEventListener("change", (ele) => {
      if (ele.target.type != "checkbox") return;
      const isChecked = ele.target.checked;
      selectAll(isChecked);
    });
  document.getElementById("pagination-next").addEventListener("click", () => {
    dbControl.renderNextRow();
  });
  document.getElementById("pagination-prev").addEventListener("click", () => {
    dbControl.renderPrevRow();
  });
  document.getElementById("pagination-list").addEventListener("click", (e) => {
    const card = e.target;
    if (!card.dataset.page || dbControl.currentPage === card.dataset.page) {
      return;
    }
    dbControl.renderPage(Number(card.dataset.page));
  });
}

function parseDataForm(data, tableName) {
  switch (tableName) {
    case "products": {
      data.price = parseInt(data.price, 10);
      data.category_id = parseInt(data.category_id, 10);
      data.discount = data.discount ? parseInt(data.discount, 10) : 0;
      break;
    }
  }

  return data;
}

function selectAll(state) {
  const inputList = document.querySelectorAll(".input-cell input");
  inputList.forEach((ele) => {
    ele.checked = state;
  });
}

export function getPageArray(totalPages, currentPage) {
  const pages = [];
  const maxVisible = 5;

  // State A: If total pages is small, just show all of them (no dots needed)
  if (totalPages <= maxVisible) {
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i);
    }
  } else {
    // State B: Near the start
    if (currentPage <= 3) {
      pages.push(1, 2, 3, 4, "...", totalPages);
    }
    // State C: Near the end
    else if (currentPage >= totalPages - 2) {
      pages.push(
        1,
        "...",
        totalPages - 3,
        totalPages - 2,
        totalPages - 1,
        totalPages,
      );
    }
    // State D: In the middle (sliding window!)
    else {
      pages.push(
        1,
        "...",
        currentPage - 1,
        currentPage,
        currentPage + 1,
        "...",
        totalPages,
      );
    }
  }
  return pages;
}

export function getSelectedIDs() {
  const checkedBoxes = document.querySelectorAll(".row-checkbox:checked");

  const ids = [];
  checkedBoxes.forEach((cb) => {
    ids.push(Number(cb.dataset.id));
  });

  return ids;
}
