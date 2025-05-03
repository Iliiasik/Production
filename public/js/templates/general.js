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

