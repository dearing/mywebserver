// ws.js is basic websocket

const wsUrl = window.location.href.replace('http', 'ws').replace('websockets.html', 'ws');
console.log(wsUrl);

ws = new WebSocket(wsUrl);

ws.onopen = () => {
  console.log('Connected');
}

// T{Message = ""}
ws.onmessage = (event) => {
  console.log('Received:', event.data);
  const p = document.createElement('p');

  info = JSON.parse(event.data)

  p.textContent = `${Date.now()}  ${info.message}`;
  document.querySelector('.events').appendChild(p);
};

ws.onerror = (error) => {
  console.error('Error:', error);
  connect();
};

ws.onclose = () => {
  console.log('Disconnected');
};

function sendMessage(message) {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(message);
  } else {
    console.log("WebSocket is not open. Cannot send message.");
  }
}

// close websocket on leaving the page
// TODO: in theory this works but not seeing the context closing out on the backend
window.onbeforeunload = () => {
  ws.close();
}
