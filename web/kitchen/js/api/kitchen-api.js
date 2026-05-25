import {
  addOrderToQueue,
  removeOrderFromQueue,
  renderOrderTotals,
} from "../ui/order_status.ui.js";

export function initSSEListen() {
  const source = new EventSource("/api/kitchen/stream");
  source.onmessage = (event) => {
    const eventData = JSON.parse(event.data);
    console.log(eventData);
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
