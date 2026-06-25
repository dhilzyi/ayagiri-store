import { fetchDatabase, fetchStatic } from "../api/kitchen-api.js";
import { prettyTimestamp } from "../ui/helpers.js";
import { getPageArray } from "./database.events.js";
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
		<td class="input-cell"><input type="checkbox" class="row-checkbox" data-id="${p.id}" /></td>
		<td>${lastRowNumber}</td>
		<td class="product-id">${p.id}</td>
		<td class="product-name">${p.name}</td>
		<td class="price">${p.price}</td>
		<td class="discount">${p.discount}</td>
		<td class="category">${dbControl.getCategoryNameByID(p.category_id)}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="updated-at">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      case "orders":
        row = `
	<tr>
		<td class="input-cell"><input type="checkbox" class="row-checkbox" data-id="${p.id}" /></td>
		<td>${lastRowNumber}</td>
		<td class="order-id">${p.id}</td>
		<td class="table-id">${p.table_id}</td>
		<td class="order-complete">${p.order_complete ? "完了" : "未完了"}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="updated-at">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      case "order_items": {
        row = `
	<tr>
		<td class="input-cell"><input type="checkbox" class="row-checkbox" data-id="${p.id}" /></td>
		<td>${lastRowNumber}</td>
		<td class="order-item-id">${p.id}</td>
		<td class="order-id">${p.order_id}</td>
		<td class="product-name">${p.product_name}</td>
		<td class="quantity">${p.quantity}</td>
		<td class="order-complete">${p.order_complete ? "完了" : "未完了"}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="updated-at">${prettyTimestamp(p.updated_at)}</td>
	</tr>
`;
        break;
      }
      case "categories": {
        row = `
	<tr>
		<td class="input-cell"><input type="checkbox" class="row-checkbox" data-id="${p.id}" /></td>
		<td>${lastRowNumber}</td>
		<td class="category_id">${p.id}</td>
		<td class="category-name">${p.name}</td>
		<td class="cateogry-english-name">${p.english_name}</td>
		<td class="created-at">${prettyTimestamp(p.created_at)}</td>
		<td class="updated-at">${prettyTimestamp(p.updated_at)}</td>
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

// TODO: add reset initial for results
export function renderResults(data, result) {
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

export async function fillPopup(product) {
  document.querySelector('input[name="name"]').value = product.name;
  document.querySelector('input[name="price"]').value = product.price;
  document.querySelector('select[name="category_id"]').value =
    product.category_id;
  document.querySelector('input[name="discount"]').value = product.discount;

  const form = document.getElementById("popup-form");
  form.dataset.action = "update";
  form.dataset.id = product.id;
}

// TODO: Move templatecache to state or somewhere
// Make it support multiple templates
let templateCache = new Map();
export async function loadAllTemplates() {
  try {
    const htmlText = await fetchStatic("templates/forms.html");

    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = htmlText;

    // Extract all <template> elements and cache them
    const templates = tempDiv.querySelectorAll("template");

    templates.forEach((template) => {
      templateCache.set(template.id, template.content);
    });

    console.log(`Loaded ${templates.length} form templates successfully`);
  } catch (error) {
    console.error("Failed to load form templates:", error);
  }
}

export async function loadCategories() {
  const tableName = "categories";
  try {
    const data = await fetchDatabase(tableName);
    dbControl.addDatabase(tableName, data);
  } catch (e) {
    console.error("Failed to load categories", e);
  }
}

export function populateCategorySelect() {
  const selectElement = document.querySelector('select[name="category_id"]');
  if (!selectElement) {
    console.warn("select element does not exist");
    return;
  }

  selectElement.innerHTML =
    '<option value="" disabled selected>カテゴリーを選択してください</option>';

  const tableData = dbControl.getTableByName("categories");

  tableData.forEach((cat) => {
    const opt = document.createElement("option");
    opt.value = cat.id;
    opt.textContent = cat.name;
    selectElement.appendChild(opt);
  });
}

export async function swapFormInput(tableName) {
  const container = document.querySelector("div.input-container");
  if (!container) {
    console.warn("container element does not exist");
    return;
  }

  container.innerHTML = "";

  const templateContent = templateCache.get(`form-${tableName}`);

  if (!templateContent) {
    container.innerHTML = `<p style="color:red">Template not found</p>`;
    return;
  }

  const clone = templateContent.cloneNode(true);
  container.appendChild(clone);

  if (tableName === "products") {
    populateCategorySelect();
  }
}
