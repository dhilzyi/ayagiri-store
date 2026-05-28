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
    console.log(JSON.stringify(eventData.payload));
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
  let results;
  switch (tableName) {
    case "products": {
      const response = await fetch("/api/admin/products");
      if (!response.ok) {
        throw Error(`response bad status: ${response.status}`);
      }
      results = response.json();
      break;
    }
    case "orders": {
      const response = await fetch("/api/admin/orders");
      if (!response.ok) {
        throw Error(`response bad status: ${response.status}`);
      }
      results = response.json();
      break;
    }
  }

  return results;
}
