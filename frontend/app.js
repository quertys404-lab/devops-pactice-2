async function loadItems() {
  const status = document.getElementById('status');
  const list = document.getElementById('items');
  status.textContent = 'Загружаем…';
  list.replaceChildren();
  try {
    const response = await fetch('/api/items');
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const items = await response.json();
    for (const item of items) {
      const li = document.createElement('li');
      li.textContent = `${item.id}. ${item.name}`;
      list.appendChild(li);
    }
    status.textContent = `Подключение работает · записей: ${items.length}`;
  } catch (error) {
    status.textContent = `Не удалось загрузить данные: ${error.message}`;
  }
}
document.getElementById('reload').addEventListener('click', loadItems);
loadItems();
