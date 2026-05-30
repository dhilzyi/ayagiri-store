import { sendNewRows } from "../api/kitchen-api.js";
import { prettyTimestamp } from "../ui/helpers.js";
import { dbControl } from "./database.js";

export function updateHeader(tableName) {
  const trHead = document.querySelector(".database-table thead tr");
  let newHead;
  switch (tableName) {
    case "products": {
      newHead = `
	<th class="input-cell" id="all-select">
		<input type="checkbox" />
	</th>
	<th>#</th>
	<th class="product-id">商品_ID <span>⇅</span></th>
	<th class="product-name">商品名</th>
	<th class="price">値段⇅</th>
	<th class="discount">割引⇅</th>
	<th class="category">カテゴリー⇅</th>
	<th class="created-at">作成日⇅</th>
	<th class="updated-at">最終変更日⇅</th>
`;
      break;
    }
    case "orders": {
      newHead = `
	<th class="input-cell" id="all-select">
		<input type="checkbox" />
	</th>
	<th>#</th>
	<th class="order-id">注文_ID</th>
	<th class="table-id">席番</th>
	<th class="order-complete">注文状況⇅</th>
	<th class="created-at">作成日⇅</th>
	<th class="updated-at">最終変更日⇅</th>
`;
      break;
    }
    case "order_items": {
      newHead = `
	<th class="input-cell" id="all-select">
		<input type="checkbox" />
	</th>
	<th>#</th>
	<th class="order-item-id">注文品_ID <span>⇅</span></th>
	<th class="order-id">注文_ID</th>
	<th class="product-name">商品名</th>
	<th class="quantity">個数</th>
	<th class="order-complete">注文状況⇅</th>
	<th class="created-at">作成日⇅</th>
	<th class="updated-at">最終変更日⇅</th>
`;
      break;
    }
    case "categories": {
      newHead = `
	<th class="input-cell" id="all-select">
		<input type="checkbox" />
	</th>
	<th>#</th>
	<th class="category-id">カテゴリー_ID <span>⇅</span></th>
	<th class="category-name">名前</th>
	<th class="cateogry-english-name">英語の名前</th>
	<th class="created-at">作成日⇅</th>
	<th class="updated-at">最終変更日⇅</th>
`;
      break;
    }
  }
  trHead.innerHTML = newHead;
}

export function addRows(data, lastRowNumber, tableName) {
  const tableBody = document.querySelector(".database-table tbody");
  const newRows = [];
  data.forEach((p) => {
    let row;
    switch (tableName) {
      case "products":
        row = `
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
        break;
      case "orders":
        row = `
	<tr>
		<th class="input-cell"><input type="checkbox" /></th>
		<td>${lastRowNumber}</td>
		<td class="order-id">${p.id}</td>
		<td class="table-id">${p.table_id}</td>
		<td class="order-complete">${p.order_complete ? "完了" : "未完了"}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="last-updated">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      case "order_items": {
        row = `
	<tr>
		<th class="input-cell"><input type="checkbox" /></th>
		<td>${lastRowNumber}</td>
		<td class="order-item-id">${p.id}</td>
		<td class="order-id">${p.order_id}</td>
		<td class="product-name">${p.product_name}</td>
		<td class="quantity">${p.quantity}</td>
		<td class="order-complete">${p.order_complete ? "完了" : "未完了"}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="last-updated">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      }
      case "categories": {
        row = `
	<tr>
		<td class="input-cell"><input type="checkbox" /></td>
		<td>${lastRowNumber}</td>
		<td class="category_id">${p.id}</td>
		<td class="category-name">${p.name}</td>
		<td class="cateogry-english-name">${p.english_name}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="last-updated">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      }
    }
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
  listContainer.innerHTML = "";

  const pageItems = getPageArray(totalPages, currentPage);

  pageItems.forEach((item) => {
    const li = document.createElement("li");

    if (item === "...") {
      // Draw dots
      li.innerHTML = `<button class="pagination-dots">...</button>`;
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

export function renderInformation(total, start, end) {
  const infoSpan = document.querySelector(
    ".database-pagination-header .left span",
  );
  infoSpan.textContent = `${total}件の結果のうち、${start}～${end}件を表示しています`;
}

export function initBtnListen() {
  let select = false;
  document
    .querySelector(".database-table thead")
    .addEventListener("change", (ele) => {
      if (ele.target.type != "checkbox") return;
      select = !select;
      selectAll(select);
      console.log("select");
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

export function initPopup() {
  const form = document.getElementById("popup-form");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const table = form.dataset.table;
    const formData = new FormData(form);

    let data = Object.fromEntries(formData);
    data = parseDataForm(data, table);
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
    // form.reset();
  });

  const modal = document.getElementById("create-row");
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.close();
    }
  });

  // DEBUG POP UP
  modal.showModal();

  document.querySelector(".action .right").addEventListener("click", (e) => {
    const cardId = e.target.id;
    if (!cardId) return;
    switch (cardId) {
      case "new-order-btn": {
        modal.showModal();
        break;
      }
    }
  });
  document.getElementById("popup-cancel-btn").addEventListener("click", () => {
    modal.close();
  });
}

function parseDataForm(data, tableName) {
  switch (tableName) {
    case "products": {
      data.price = parseInt(data.price, 10);
      data.category_id = parseInt(data.category_id, 10);
      data.discount = data.discount ? parseInt(data.discount, 10) : 0;
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

function renderResults(data, result) {
  const submitResult = document.querySelector("div.submit-result");
  switch (result) {
    case "success": {
      submitResult.innerHTML = `
	<h4 class="result-title">送信が完了しました</h4>
	<p class="result-code">ステータスコード: <span>${data.status}</span></p>
	<p class="result-help">データは正常に登録されました。</p>
`;
      submitResult.classList.remove("error");
      submitResult.classList.add("success");
      break;
    }
    case "error": {
      submitResult.innerHTML = `
	<h4 class="result-title">送信に失敗しました</h4>
	<p class="result-code">エラーコード: <span>${data.status}</span></p>
	<p class="result-details">エラー詳細: ${data.error}</p>
	<p class="result-help">
		もう一度お試しください。<br />
		解決しない場合は、システム管理者までご連絡ください。
	</p>
`;
      submitResult.classList.remove("success");
      submitResult.classList.add("error");
      break;
    }
  }
}

export async function initSelect() {
  const databaseSelect = document.getElementById("database-select");
  const perPage = document.getElementById("result-per-page-opt");
  databaseSelect.addEventListener("change", () => {
    console.log(databaseSelect.value);
    dbControl.changeTable(databaseSelect.value);
  });
  perPage.addEventListener("change", async () => {
    await dbControl.setItemsPerPage(Number(perPage.value));
    dbControl.renderPage(dbControl.currentPage);
  });
  await dbControl.setItemsPerPage(Number(perPage.value));
  console.log(perPage.value);
  await dbControl.changeTable(databaseSelect.value);
  dbControl.renderPage(1);
}
