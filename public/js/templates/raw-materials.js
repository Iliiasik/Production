document.addEventListener("DOMContentLoaded", function () {

    // УДАЛЕНИЕ

    document.querySelectorAll(".delete-btn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            let rawmaterialId = this.getAttribute("data-id");

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
                    fetch(`/raw-materials/delete/${rawmaterialId}`, {
                        method: "DELETE"
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                // Удаляем строку только после успешного удаления на сервере (чтобы не было ошибки окон)
                                const row = document.getElementById(`row-${rawmaterialId}`);
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
            let rawMaterialId = this.getAttribute("data-id");
            console.log(`Edit button clicked for rawMaterialId: ${rawMaterialId}`);

            // Загружаем данные о сырье
            console.log(`Fetching data for raw material with ID: ${rawMaterialId}`);
            fetch(`/raw-materials/get/${rawMaterialId}`)
                .then(response => response.json())
                .then(data => {
                    console.log('Raw material data received:', data);

                    // Обновляем ключи для правильного доступа к данным
                    if (data.success && data.rawmaterial) {
                        const rawMaterial = data.rawmaterial;
                        console.log('Raw material:', rawMaterial);

                        // Проверка на существование unit_id
                        console.log('Raw material unit_id:', rawMaterial.unit_id);

                        // Загружаем доступные единицы измерения для выпадающего списка
                        console.log('Fetching units list...');
                        fetch("/units/list")
                            .then(response => response.json())
                            .then(unitsData => {
                                console.log('Units data received:', unitsData);

                                if (!unitsData.success) {
                                    Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
                                    return;
                                }

                                // Создаем выпадающий список
                                let unitOptions = unitsData.units.map(unit => {
                                    const selected = unit.id === rawMaterial.unit_id ? 'selected' : '';
                                    console.log(`Creating option for unit: ${unit.name}, selected: ${selected}`);
                                    return `<option value="${unit.id}" ${selected}>${unit.name}</option>`;
                                }).join("");
                                console.log('Unit options generated:', unitOptions);

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
                                            margin: 0;
                                        }
                                        /* Стиль для ячеек таблицы */
                                        .table-cell {
                                            padding: 10px;
                                            border: 1px solid #ddd;
                                        }
                                    </style>
                                    <form id="editForm">
                                        <table style="width:100%; border-collapse: collapse; table-layout: fixed;">
                                            <!-- Первая строка: Заголовки столбцов -->
                                            <tr style="background-color: #f4f4f4; text-align: left;">
                                                <th class="table-cell">Название</th>
                                                <th class="table-cell">Единица измерения</th>
                                                                                                <th class="table-cell">Количество</th>                                       
                                                <th class="table-cell">Сумма</th>                                       
                                       
                                            </tr>
                                            <!-- Вторая строка: Поля ввода -->
                                            <tr>
                                                <td class="table-cell">
                                                    <input type="text" id="materialName" class="input-field" 
                                                        value="${rawMaterial.name}" placeholder="Введите название">
                                                </td>
                                                <td class="table-cell">
                                                    <select id="unitId" class="input-field">
                                                        ${unitOptions}
                                                    </select>
                                                </td>
                                                <td class="table-cell">
                                                <input type="number" id="quantity" class="input-field" 
                                                        value="${rawMaterial.quantity}" step="0.01">
                                                </td>
                                                <td class="table-cell">
                                                    <input type="number" id="totalAmount" class="input-field" 
                                                        value="${rawMaterial.total_amount}" step="0.01">
                                                </td>
                                            </tr>
                                        </table>
                                    </form>
                                `,
                                    showCancelButton: true,
                                    confirmButtonText: 'Сохранить изменения',
                                    cancelButtonText: 'Отмена',
                                    preConfirm: () => {
                                        const name = document.getElementById('materialName').value;
                                        const unitId = document.getElementById('unitId').value;
                                        const quantity = document.getElementById('quantity').value;
                                        const totalAmount = document.getElementById('totalAmount').value;

                                        console.log('Form data to save:', { name, unitId, quantity, totalAmount });

                                        if (!name || !unitId || !quantity || !totalAmount) {
                                            Swal.showValidationMessage('Заполните все поля');
                                            return false;
                                        }

                                        return fetch(`/raw-materials/edit/${rawMaterialId}`, {
                                            method: 'POST',
                                            headers: {
                                                'Content-Type': 'application/json'
                                            },
                                            body: JSON.stringify({
                                                name,
                                                unit_id: parseInt(unitId),
                                                quantity: parseFloat(quantity),
                                                total_amount: parseFloat(totalAmount)

                                            })
                                        })
                                            .then(response => response.json())
                                            .then(data => {
                                                console.log('Save result:', data);

                                                if (data.success) {
                                                    return data;
                                                } else {
                                                    Swal.showValidationMessage(data.error || 'Не удалось сохранить изменения');
                                                }
                                            })
                                            .catch(error => {
                                                console.error('Error during save:', error);
                                                Swal.showValidationMessage('Ошибка при сохранении');
                                            });
                                    },
                                    focusConfirm: false,
                                    customClass: {
                                        popup: 'popup-class', // Класс для модального окна
                                        confirmButton: 'custom-button', // Класс для кнопки подтверждения
                                        cancelButton: 'custom-button' // Класс для кнопки отмены
                                    },
                                    width: '900px', // Увеличена ширина модального окна
                                }).then((result) => {
                                    if (result.isConfirmed) {
                                        location.reload();
                                    }
                                });

                            }).catch(error => {
                            console.error('Error fetching units data:', error);
                            Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
                        });

                    } else {
                        console.error('Error loading raw material data:', data);
                        Swal.fire("Ошибка!", "Не удалось загрузить данные для редактирования.", "error");
                    }
                })
                .catch(error => {
                    console.error('Error fetching raw material data:', error);
                    Swal.fire("Ошибка!", "Произошла ошибка при загрузке данных.", "error");
                });

        });
    });

    // ДОБАВЛЕНИЕ

    document.getElementById('addBtn').addEventListener('click', function () {
        // Загружаем доступные единицы измерения для выпадающего списка
        fetch("/units/list")
            .then(response => response.json())
            .then(unitsData => {
                if (!unitsData.success) {
                    Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
                    return;
                }

                // Создаем выпадающий список единиц измерения
                let unitOptions = unitsData.units.map(unit => {
                    return `<option value="${unit.id}">${unit.name}</option>`;
                }).join("");

                // Открываем модальное окно для добавления записи
                Swal.fire({
                    title: 'Добавить запись',
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
                        /* Стиль для контейнеров полей */
                        .form-group {
                            margin-bottom: 16px;
                        }
                    </style>
                    <form id="addForm">
                        <!-- Поле "Название" -->
                        <div class="form-group">
                            <label for="materialName" style="display:block; margin-bottom: 5px; font-weight: bold;">Название:</label>
                            <input type="text" id="materialName" class="input-field" placeholder="Введите название">
                        </div>

                        <!-- Поле "Единица измерения" -->
                        <div class="form-group">
                            <label for="unitId" style="display:block; margin-bottom: 5px; font-weight: bold;">Единица измерения:</label>
                            <select id="unitId" class="input-field">
                                ${unitOptions}
                            </select>
                        </div>                        
                    </form>
                `,
                    showCancelButton: true,
                    confirmButtonText: 'Добавить',
                    cancelButtonText: 'Отмена',
                    preConfirm: () => {
                        // Получаем данные формы
                        const name = document.getElementById('materialName').value;
                        const unitId = document.getElementById('unitId').value;


                        // Проверяем, чтобы все поля были заполнены
                        if (!name || !unitId) {
                            Swal.showValidationMessage('Пожалуйста, заполните все поля');
                            return false;
                        }

                        // Отправляем данные на сервер
                        return fetch('/raw-materials/add', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify({
                                name,
                                unit_id: parseInt(unitId),
                            })
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
                    width: '600px', // Установлена оптимальная ширина модального окна
                }).then((result) => {
                    if (result.isConfirmed) {
                        location.reload();
                    }
                });

            }).catch(error => {
            console.error('Error fetching units data:', error);
            Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
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


