document.addEventListener("DOMContentLoaded", function () {

    // УДАЛЕНИЕ

    document.querySelectorAll(".delete-btn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            let unitId = this.getAttribute("data-id");

            Swal.fire({
                title: "Вы уверены?",
                text: "Это действие нельзя отменить!",
                icon: "warning",
                showCancelButton: true,
                confirmButtonColor: "#d33",
                cancelButtonColor: "#3085d6",
                confirmButtonText: "Да, удалить!",
                cancelButtonText: "Отмена",
                customClass: {
                    confirmButton: 'custom-button',
                    cancelButton: 'custom-button'
                }
            }).then((result) => {
                if (result.isConfirmed) {
                    fetch(`/units/delete/${unitId}`, {
                        method: "DELETE"
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                // Удаляем строку только после успешного удаления на сервере (чтобы не было ошибки окон)
                                const row = document.getElementById(`row-${unitId}`);
                                if (row) {
                                    row.remove();
                                }
                                Swal.fire({
                                    title: "Удалено!",
                                    text: "Запись успешно удалена.",
                                    icon: "success",
                                    timer: 1000,
                                    timerProgressBar: true,
                                    toast: true,
                                    position: "top-end",
                                    showConfirmButton: false
                                });

                                setTimeout(function() {
                                    location.reload();
                                }, 1000);
                            } else {
                                Swal.fire("Ошибка!", data.error || "Не удалось удалить запись.", "error");
                            }
                        })
                        .catch(error => {
                            console.error("Error:", error);
                            Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
                        });
                }
            });
        });
    });

    // РЕДАКТИРОВАНИЕ

    document.querySelectorAll(".edit-btn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            let unitId = this.getAttribute("data-id");

            fetch(`/units/get/${unitId}`)
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        const unit = data.unit;

                        // Открываем модальное окно с данными для редактирования
                        Swal.fire({
                            title: 'Редактировать запись',
                            html: `
                            <style>
                                /* Общий стиль для всех полей ввода */
                                .input-field {
                                    width: 100%;
                                    box-sizing: border-box;
                                    padding: 8px;
                                    border: 1px solid #ccc;
                                    border-radius: 4px;
                                    font-size: 14px;
                                }
                                /* Стиль для ячеек таблицы */
                                .table-cell {
                                    padding: 10px;
                                    border: 1px solid #ddd;
                                }
                            </style>
                            <form id="editForm">
                                <table style="width:100%; border-collapse: collapse; table-layout: fixed;">
                                    <tr style="background-color: #f4f4f4; text-align: left;">
                                        <th class="table-cell">Название</th>
                                    </tr>
                                    <tr>
                                        <td class="table-cell">
                                            <input type="text" id="unitName" class="input-field" 
                                                value="${unit.name}" placeholder="Введите название">
                                        </td>
                                    </tr>
                                </table>
                            </form>
                        `,
                            showCancelButton: true,
                            confirmButtonText: 'Сохранить изменения',
                            cancelButtonText: 'Отмена',
                            preConfirm: () => {
                                // Получаем обновленные данные формы
                                const name = document.getElementById('unitName').value;

                                if (!name) {
                                    Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
                                    return false;
                                }

                                // Отправляем обновленные данные на сервер
                                return fetch(`/units/edit/${unitId}`, {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json'
                                    },
                                    body: JSON.stringify({ name })
                                })
                                    .then(response => response.json())
                                    .then(data => {
                                        if (data.success) {
                                            return data;
                                        } else {
                                            Swal.showValidationMessage(data.error || 'Не удалось сохранить изменения');
                                        }
                                    })
                                    .catch(error => {
                                        Swal.showValidationMessage('Произошла ошибка при сохранении');
                                    });
                            },
                            focusConfirm: false,
                            customClass: {
                                popup: 'popup-class', // Класс для модального окна
                                confirmButton: 'custom-button', // Класс для кнопки подтверждения
                                cancelButton: 'custom-button' // Класс для кнопки отмены
                            },
                            width: '600px', // Установлена оптимальная ширина модального окна
                        }).then((result) => {
                            if (result.isConfirmed) {
                                location.reload();
                            }
                        });
                    } else {
                        Swal.fire("Ошибка!", "Не удалось загрузить данные для редактирования.", "error");
                    }
                })
                .catch(error => {
                    console.error("Error:", error);
                    Swal.fire("Ошибка!", "Произошла ошибка при загрузке данных для редактирования.", "error");
                });
        });
    });

    // ДОБАВЛЕНИЕ

    document.getElementById('addBtn').addEventListener('click', function () {

        Swal.fire({
            title: 'Добавить запись',
            html: `
            <form id="addForm">
                <div>
                    <label for="unitName" style="display:block; margin-bottom: 5px; font-weight: bold;">Название:</label>
                    <input type="text" id="unitName" class="swal2-input" placeholder="Введите название" style="padding: 8px; width: 200px; font-size: 14px;">
                </div>
            </form>
        `,
            showCancelButton: true,
            confirmButtonText: 'Добавить',
            cancelButtonText: 'Отмена',
            preConfirm: () => {
                // Получаем данные формы
                const name = document.getElementById('unitName').value;

                // Проверяем, чтобы поле не было пустым
                if (!name) {
                    Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
                    return false;
                }

                // Отправляем данные на сервер
                return fetch('/units/add', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            return data;
                        } else {
                            Swal.showValidationMessage(data.error || 'Не удалось добавить запись');
                        }
                    })
                    .catch(error => {
                        Swal.showValidationMessage('Произошла ошибка при добавлении');
                    });
            },
            focusConfirm: false,
            customClass: {
                popup: 'popup-class',
                confirmButton: 'custom-button',
                cancelButton: 'custom-button'
            },
        }).then((result) => {
            if (result.isConfirmed) {
                location.reload();
            }
        });
    });

    const style = document.createElement('style');
    style.innerHTML = `
    .swal2-popup {
        width: 500px; 
        padding: 20px; 
        font-size: 16px; 
    }
    .swal2-input {
        border: 1px solid #ccc; 
        border-radius: 4px; 
        padding: 8px; 
        width: 100%; 
    }
    .popup-class {
        font-size: 16px;
    }
    .custom-button {
        font-family: 'TildaSans', sans-serif !important; /* Применяем шрифт к кнопкам */
        font-size: 16px !important; /* Устанавливаем размер шрифта для кнопок */
    }
`;
    document.head.appendChild(style);

});


