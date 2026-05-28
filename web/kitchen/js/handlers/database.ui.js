import { prettyTimestamp } from "../ui/helpers.js";
import { dbControl } from "./database.js";

export function addRows(products, lastRowNumber) {
  const tableBody = document.querySelector(".database-table tbody");
  const newRows = [];
  products.forEach((p) => {
    const row = `
	<tr>
		<th class="input-cell"><input type="checkbox" /></th>
		<td>${lastRowNumber}</td>
		<td class="product-id">${p.id}</td>
		<td class="product-name">${p.name}</td>
		<td>${p.price}</td>
		<td class="discount">${p.discount}</td>
		<td class="category">${p.category_name}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="last-updated">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
    newRows.push(row);
    lastRowNumber++;
  });
  tableBody.innerHTML = newRows.join("");
}

function getPageArray(totalPages, currentPage) {
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

export function renderPaginationBar(totalPages, currentPage) {
  const listContainer = document.getElementById("pagination-list");
  listContainer.innerHTML = ""; // Clear the old list

  const pageItems = getPageArray(totalPages, currentPage);

  pageItems.forEach((item) => {
    const li = document.createElement("li");

    if (item === "...") {
      // Draw dots
      li.innerHTML = `<span class="pagination-dots">...</span>`;
    } else {
      // Draw standard page button
      const isActive = item === currentPage ? "active" : "";
      li.innerHTML = `<button class="page-btn ${isActive}" data-page="${item}">${item}</button>`;
    }

    listContainer.appendChild(li);
  });

  // Toggle Prev/Next buttons disabled state based on page boundaries
  document.getElementById("pagination-prev").disabled = currentPage === 1;
  document.getElementById("pagination-next").disabled =
    currentPage === totalPages;
}

export function btnListen() {
  let select = false;
  document.getElementById("all-select").addEventListener("change", () => {
    select = !select;
    selectAll(select);
  });
  document.getElementById("pagination-next").addEventListener("click", () => {
    dbControl.renderNextRow("products");
  });
  document.getElementById("pagination-prev").addEventListener("click", () => {
    dbControl.renderPrevRow("products");
  });
  document.getElementById("pagination-list").addEventListener("click", (e) => {
    const card = e.target;
    if (!card.dataset.page || dbControl.currentPage === card.dataset.page) {
      return;
    }
    dbControl.renderPage("products", Number(card.dataset.page));
  });
}

function selectAll(state) {
  const inputList = document.querySelectorAll(".input-cell input");
  inputList.forEach((ele) => {
    ele.checked = state;
  });
}
