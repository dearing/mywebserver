const eventSource = new EventSource('/events');

eventSource.onmessage = (event) => {
  console.log('Received:', event);
  const p = document.createElement('p');
  p.textContent = `${Date.now()} => ${event.data}`;
  document.querySelector('.events').appendChild(p);
};

eventSource.addEventListener('refresh', (event) => {
  console.log('Refresh event received:', event.data);
  window.location.reload();
});

eventSource.onerror = (error) => {
  console.error('EventSource failed:', error);
  window.location.reload();
};