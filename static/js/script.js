console.log("Frontend script loaded!");

document.addEventListener('DOMContentLoaded', () => {
    const routeSelect = document.getElementById('route'); // Новый выпадающий список
    const getRecommendationsButton = document.getElementById('getRecommendations');
    const stopsContainer = document.querySelector('.stops');

    // Функция для загрузки номеров маршрутов
    async function loadRouteNumbers() {
        try {
            const response = await fetch('/getRouteNumbers');
            const data = await response.json();

            // Очистка списка перед добавлением новых элементов
            routeSelect.innerHTML = '';

            console.log(111, data);
            // Добавление номеров маршрутов в выпадающий список
            data.routeNumbers.forEach(route => {
                const option = document.createElement('option');
                option.value = route;
                option.textContent = route;
                routeSelect.appendChild(option);
            });
        } catch (error) {
            console.error('Ошибка при загрузке номеров маршрутов:', error);
        }
    }

    // Функция для обновления остановок
    async function updateStops() {
        try {
            // Отправка запроса на бэкенд
            const response = await fetch('/getStops', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ routeNumber: "1", timePeriod: "now" }), // Пример данных
            });

            const data = await response.json();

            // Очистка предыдущих остановок
            stopsContainer.innerHTML = '';

            // Отображение остановок
            data.stops.forEach((stop, index) => {
                const stopElement = document.createElement('div');
                stopElement.className = 'stop';

                // Определение цвета точки в зависимости от загруженности
                if (stop.Workload < 60) {
                    stopElement.classList.add('green');
                } else if (stop.Workload >= 60 && stop.Workload <= 80) {
                    stopElement.classList.add('yellow');
                } else {
                    stopElement.classList.add('red');
                }

                // Позиционирование точки
                stopElement.style.left = `${(index / (data.stops.length - 1)) * 100}%`;
                stopsContainer.appendChild(stopElement);
            });
        } catch (error) {
            console.error('Ошибка при обновлении остановок:', error);
        }
    }

    loadRouteNumbers();
    // Первоначальное обновление
    updateStops();

    // Автоматическое обновление каждые 5 минут (300000 миллисекунд)
    setInterval(updateStops, 300000);
});


// document.addEventListener('DOMContentLoaded', () => {
//     const getRecommendationsButton = document.getElementById('getRecommendations');
//     const stopsContainer = document.querySelector('.stops');
//
//     // Обработчик кнопки "Получить рекомендации"
//     getRecommendationsButton.addEventListener('click', async () => {
//         const routeNumber = document.getElementById('routeNumber').value;
//         const timePeriod = document.getElementById('timePeriod').value;
//
//         // Отправка запроса на бэкенд
//         const response = await fetch('/getStops', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({ routeNumber, timePeriod }),
//         });
//
//         const data = await response.json();
//
//         // Очистка предыдущих остановок
//         stopsContainer.innerHTML = '';
//
//         // Отображение остановок
//         data.stops.forEach((stop, index) => {
//             const stopElement = document.createElement('div');
//             stopElement.className = 'stop';
//
//             // Определение цвета точки в зависимости от загруженности
//             if (stop.Workload < 60) {
//                 stopElement.classList.add('green');
//             } else if (stop.Workload >= 60 && stop.Workload < 80) {
//                 stopElement.classList.add('yellow');
//             } else {
//                 stopElement.classList.add('red');
//             }
//
//             // Позиционирование точки
//             stopElement.style.left = `${(index / (data.stops.length - 1)) * 100}%`;
//             stopsContainer.appendChild(stopElement);
//         });
//     });
// });
