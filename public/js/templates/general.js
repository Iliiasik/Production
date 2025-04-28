// Функция для загрузки данных с сервера
function loadData(url) {
    return fetch(url)
        .then(response => response.json())
        .then(data => {
            if (!data.success) {
                Swal.fire("Ошибка!", "Не удалось загрузить данные.", "error");
                return null;
            }
            return data;
        })
        .catch(error => {
            console.error("Error:", error);
            Swal.fire("Ошибка!", "Произошла ошибка при загрузке данных.", "error");
            return null;
        });
}

// Функция для отображения модального окна с формой
function showModal(title, formHtml, preConfirm, width = '600px') {
    return Swal.fire({
        title,
        html: formHtml,
        showCancelButton: true,
        confirmButtonText: title === 'Добавить запись' ? 'Добавить' : 'Сохранить изменения',
        cancelButtonText: 'Отмена',
        preConfirm,
        customClass: {
            popup: 'popup-class',
            confirmButton: 'custom-button',
            cancelButton: 'custom-button'
        },
        width,
    });
}

function secureFetch(url, options = {}) {
    return fetch(url, {
        credentials: 'include', // важно: чтобы cookie (в т.ч. JWT) передавались
        ...options,
    })
        .then(async response => {
            if (!response.ok) {
                let errorData = {};
                try {
                    errorData = await response.json();
                } catch (_) {}

                // Обработка ошибок авторизации и доступа
                if (response.status === 401) {
                    Swal.fire({
                        title: "Не авторизованы",
                        text: errorData.error || "Пожалуйста, войдите в систему",
                        icon: "warning",
                        confirmButtonText: "Ок",
                        customClass: {
                            popup: 'popup-class',
                            confirmButton: 'custom-button',
                            cancelButton: 'custom-button'
                        },
                    });
                } else if (response.status === 403) {
                    Swal.fire({
                        title: "Доступ запрещён",
                        text: errorData.error || "У вас нет прав для выполнения действия",
                        icon: "error",
                        confirmButtonText: "Ок",
                        customClass: {
                            popup: 'popup-class',
                            confirmButton: 'custom-button',
                            cancelButton: 'custom-button'
                        },
                    });
                } else {
                    Swal.fire({
                        title: "Ошибка",
                        text: errorData.error || "Произошла неизвестная ошибка",
                        icon: "error",
                        confirmButtonText: "Ок",
                        customClass: {
                            popup: 'popup-class',
                            confirmButton: 'custom-button',
                            cancelButton: 'custom-button'
                        },
                    });
                }

                throw errorData;
            }

            // Успешный ответ (json или text — можно сделать по ситуации)
            const contentType = response.headers.get("content-type");
            if (contentType && contentType.includes("application/json")) {
                return response.json();
            } else {
                return response.text();
            }
        })
        .catch(error => {
            console.error("SecureFetch Error:", error);
            throw error; // опционально: можешь возвращать null, если не хочешь падать дальше
        });
}

document.querySelectorAll('a[href^="/"]').forEach(link => {
    link.addEventListener('click', (event) => {
        event.preventDefault(); // отменяем стандартное поведение ссылки

        const url = link.getAttribute('href');

        // Выполняем запрос с проверкой прав доступа
        secureFetch(url)
            .then(data => {
                // Если данные получены, переходим по маршруту
                if (data) {
                    window.location.href = url; // переходим по URL
                }
            })
    });
});
