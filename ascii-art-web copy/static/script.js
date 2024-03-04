// Отправляем GET-запрос при загрузке страницы
sendGetRequest();

// Функция для отправки GET-запроса и обработки ответа
function sendGetRequest() {
    fetch('/ascii-art') // Отправляем GET-запрос на URL '/ascii-art'
        .then(response => response.text()) // Получаем ответ и преобразуем его в текст
        .then(data => { // Обрабатываем полученные данные
            document.getElementById('output').innerText = data; // Выводим данные на страницу
        })
        .catch(error => console.error('Error fetching data:', error)); // В случае ошибки выводим сообщение в консоль
}
function download() {
    fetch('/ascii-art') // Отправляем GET-запрос на URL '/ascii-art'
        .then(response => response.text()) // Получаем ответ и преобразуем его в текст
        .then(data => { // Обрабатываем полученные данные
            // Создаем новый Blob из текстовых данных
            const blob = new Blob([data], { type: 'text/plain' });
            // Создаем ссылку на Blob
            const url = window.URL.createObjectURL(blob);
            // Создаем ссылку для скачивания файла
            const a = document.createElement('a');
            a.href = url;
            a.download = 'ascii-art.txt'; // Устанавливаем имя файла
            // Добавляем ссылку в документ
            document.body.appendChild(a);
            // Кликаем по ссылке для скачивания файла
            a.click();
            // Удаляем ссылку из документа
            window.URL.revokeObjectURL(url);
            // Удаляем элемент ссылки из документа
            a.remove();
        })
        .catch(error => console.error('Error fetching data:', error));
}