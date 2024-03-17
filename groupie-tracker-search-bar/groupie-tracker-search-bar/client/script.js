

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

window.onload = function() {
    var map = L.map('map').setView([48.0196, 66.9237], 2);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    var dataElement = document.getElementById("data");
    var cityString = dataElement.textContent.trim(); // Получаем текст из элемента и удаляем лишние пробелы
    var cityNames = cityString.split(" "); // Разбиваем строку на отдельные названия городов
    for (var i = 0; i < cityNames.length; i++) {
        cityNames[i] = cityNames[i].replace(/-/g, ",");
    }
    function geocodeCity(city) {
        var url = "https://nominatim.openstreetmap.org/search?q=" + encodeURIComponent(city) + "&format=json&limit=1";
        console.log(url);
        fetch(url)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                if (data.length > 0) {
                    var lat = parseFloat(data[0].lat);
                    var lon = parseFloat(data[0].lon);
                    L.marker([lat, lon]).addTo(map)
                        .bindPopup(city); 
                } else {
                    console.error("Геокодирование не удалось для города:", city);
                }
            })
            .catch(error => {
                console.error("Произошла ошибка при геокодировании:", error);
            });
    }

    cityNames.forEach(function (city) {
        geocodeCity(city);
    });
}


