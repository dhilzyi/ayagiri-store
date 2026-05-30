import {
  addOrderToQueue,
  removeOrderFromQueue,
  renderOrderTotals,
} from "../ui/order_status.ui.js";

export function initSSEListen() {
  const eventSource = new EventSource("/api/kitchen/stream");
  const statusContainer = document.getElementById("connection-status");
  const statusText = statusContainer.querySelector(".status-text");

  eventSource.onopen = function () {
    statusContainer.className = "status-connected";
    statusText.textContent = "接続中 (Online)";
    console.log("connected");
  };
  eventSource.onerror = function () {
    statusContainer.className = "status-disconnected";
    statusText.textContent = "接続切れ (Offline)";
    console.log("disconnected");
    // Note: EventSource will automatically try to reconnect in the background.
  };

  eventSource.onmessage = (event) => {
    const eventData = JSON.parse(event.data);
    switch (eventData.type) {
      case "ADD_ORDER":
        console.log("ADD", eventData.payload.order_id);
        addOrderToQueue(eventData.payload);
        renderOrderTotals();
        break;
      case "COMPLETE_ORDER":
        console.log("COMPLETE", eventData.payload.order_id);
        removeOrderFromQueue(eventData.payload.order_id);
        break;
      default:
        break;
    }
  };
}

export async function sendComplete(orderID) {
  const res = await fetch(`/api/complete_orders?order_id=${orderID}`, {
    method: "POST",
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(
      errorData.message || `Server returned status ${res.status}`,
    );
  }
  console.log(res.status);
}

export async function fetchDatabase(tableName) {
  let url;
  switch (tableName) {
    case "products": {
      url = "/api/admin/products";
      break;
    }
    case "orders": {
      url = "/api/admin/orders";
      break;
    }
    case "order_items": {
      url = "/api/admin/order_items";
      break;
    }
    case "categories": {
      url = "/api/admin/categories";
      break;
    }
  }

  const response = await fetch(url);
  if (!response.ok) {
    throw Error(`response bad status: ${response.status}`);
  }
  const results = response.json();

  return results;
}

export async function sendNewRows(data, tableName) {
  const url = `/api/${tableName}`;

  const response = await fetch(url, { method: "POST", body: data });
  if (!response.ok) {
    const err = new Error(`Failed to send data`);

    err.body = response.json();
    err.status = response.status;

    throw err;
  }
  return response;
}
