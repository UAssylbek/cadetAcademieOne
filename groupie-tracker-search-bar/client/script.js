

function toggleInfo(type) {
    var element = document.getElementById(type + "-det");
    if (element) {
        // Сначала скрываем все элементы
        document.querySelectorAll('.location-details').forEach(function (el) {
            el.classList.remove("d-block");
            el.classList.add("d-none");
        });

        // Затем открываем только выбранный элемент
        element.classList.remove("d-none");
        element.classList.add("d-block");
    }
}

var map; // Глобальная переменная для хранения объекта карты

function initMap(locationsss) {
    map = L.map('map').setView([48.0196, 66.9237], 2); // Используем глобальную переменную для хранения объекта карты

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    // Геокодирование каждого города из данных
    locationsss.forEach(function(location) {
        geocodeCity(location.city);
    });
}

// Функция для геокодирования городов и стран
function geocodeCity(city) {
    var url = "https://nominatim.openstreetmap.org/search?q=" + encodeURIComponent(city) + "&format=json&limit=1";

    fetch(url)
        .then(response => response.json())
        .then(data => {
            if (data && data.length > 0) {
                var lat = parseFloat(data[0].lat);
                var lon = parseFloat(data[0].lon);
                L.marker([lat, lon]).addTo(map)
                    .bindPopup(city); // Используем глобальную переменную для доступа к объекту карты
            } else {
                console.error("Геокодирование не удалось для города:", city);
            }
        })
        .catch(error => {
            console.error("Произошла ошибка при геокодировании:", error);
        });
}

// Получаем данные с сервера
fetch('/').then(response => response.json())
           .then(data => initMap(data))
           .catch(error => console.error("Произошла ошибка при получении данных:", error));
