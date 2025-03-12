document.addEventListener("DOMContentLoaded", function () {
    const container = document.body; // Контейнер для делегирования событий

    // УДАЛЕНИЕ с использованием делегирования событий
    container.addEventListener('click', function (e) {
        if (e.target && e.target.closest(".delete-btn")) {
            e.preventDefault();
            let ingredientId = e.target.closest(".delete-btn").getAttribute("data-id");

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
                    fetch(`/ingredients/delete/${ingredientId}`, {
                        method: "DELETE"
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                // Удаляем строку из таблицы без перезагрузки страницы
                                const row = document.getElementById(`row-${ingredientId}`);
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

                                // Обновляем таблицу, чтобы отразить изменения
                                updateIngredientsTable();
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
        }
    });

    // РЕДАКТИРОВАНИЕ с использованием делегирования событий
    container.addEventListener('click', function (e) {
        if (e.target && e.target.closest(".edit-btn")) {
            e.preventDefault();
            let ingredientId = e.target.closest(".edit-btn").getAttribute("data-id");

            fetch(`/ingredients/get/${ingredientId}`)
                .then(response => response.json())
                .then(data => {
                    if (data.success && data.ingredient) {
                        const ingredient = data.ingredient;

                        // Загружаем список готовой продукции и сырья
                        Promise.all([
                            fetch("/raw-materials/list").then(res => res.json()),
                            fetch(`/ingredients/used-raw-materials/${ingredient.product_id}`).then(res => res.json())
                        ]).then(([rawMaterialsData, usedRawMaterialsData]) => {
                            // Проверяем валидность данных
                            if (!rawMaterialsData.success || !usedRawMaterialsData.success) {
                                Swal.fire("Ошибка!", "Не удалось загрузить необходимые данные.", "error");
                                return;
                            }

                            // Получаем массив ID используемого сырья
                            const usedRawMaterialIds = usedRawMaterialsData.used_raw_materials.map(material => material.id);

                            // Фильтруем сырье: исключаем те ингредиенты, которые уже используются (кроме текущего)
                            const availableRawMaterials = rawMaterialsData.raw_materials.filter(rawMaterial => {
                                return !usedRawMaterialIds.includes(rawMaterial.id) || rawMaterial.id === ingredient.raw_material_id;
                            });


                            // Создаем выпадающий список для сырья
                            let rawMaterialOptions = availableRawMaterials.map(rawMaterial => {
                                const selected = rawMaterial.id === ingredient.raw_material_id ? 'selected' : '';
                                return `<option value="${rawMaterial.id}" ${selected}>${rawMaterial.name}</option>`;
                            }).join("");

                            // Открываем модальное окно для редактирования
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
                                    
                                <table style="width:100%; border-collapse: collapse; table-layout: fixed;">
    <!-- Первая строка: Заголовки столбцов -->
    <tr style="background-color: #f4f4f4; text-align: left;">
        <th class="table-cell">Сырье</th>
        <th class="table-cell">Количество</th>
    </tr>
    <!-- Вторая строка: Поля ввода -->
    <tr>

        <td class="table-cell">
            <select id="rawMaterialId" class="input-field">
                ${rawMaterialOptions}
            </select>
        </td>
        <td class="table-cell">
            <input id="quantity" type="number" class="input-field" placeholder="Введите количество" value="${ingredient.quantity}">
        </td>
    </tr>
</table>
                            `,
                                showCancelButton: true,
                                confirmButtonText: 'Сохранить изменения',
                                cancelButtonText: 'Отмена',
                                preConfirm: () => {
                                    const rawMaterialId = document.getElementById('rawMaterialId').value;
                                    const quantity = document.getElementById('quantity').value;

                                    if (!rawMaterialId || !quantity) {
                                        Swal.showValidationMessage('Заполните все поля');
                                        return false;
                                    }

                                    return fetch(`/ingredients/edit/${ingredientId}`, {
                                        method: 'POST',
                                        headers: {
                                            'Content-Type': 'application/json'
                                        },
                                        body: JSON.stringify({
                                            raw_material_id: parseInt(rawMaterialId),
                                            quantity: parseFloat(quantity)
                                        })
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
                                            console.error('Error during save:', error);
                                            Swal.showValidationMessage('Ошибка при сохранении');
                                        });
                                },focusConfirm: false,
                                customClass: {
                                    popup: 'popup-class', // Класс для модального окна
                                    confirmButton: 'custom-button', // Класс для кнопки подтверждения
                                    cancelButton: 'custom-button' // Класс для кнопки отмены
                                },
                            }).then((result) => {
                                if (result.isConfirmed) {
                                    updateIngredientsTable();
                                }
                            });
                        });
                    } else {
                        Swal.fire("Ошибка!", "Не удалось загрузить данные для редактирования.", "error");
                    }
                })
                .catch(error => {
                    Swal.fire("Ошибка!", "Произошла ошибка при загрузке данных.", "error");
                });
        }
    });
    // ДОБАВЛЕНИЕ
    document.getElementById('addBtn').addEventListener('click', function () {
        console.log("Add button clicked. Fetching data for adding ingredient...");

        const productId = document.getElementById('productSelect').value;
        const productName = document.getElementById('productSelect').options[
            document.getElementById('productSelect').selectedIndex
            ].text;

        console.log(`Selected product ID: ${productId}, Name: "${productName}"`);

        // Загружаем список сырья
        fetch("/raw-materials/list")
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(rawMaterialsData => {
                console.log("Raw materials data loaded:", rawMaterialsData);

                if (!rawMaterialsData.success || !Array.isArray(rawMaterialsData.raw_materials)) {
                    console.error("Invalid raw materials data:", rawMaterialsData);
                    Swal.fire("Ошибка!", "Не удалось загрузить список сырья.", "error");
                    return;
                }

                // Загружаем список использованного сырья для текущего продукта
                fetch(`/ingredients/used-raw-materials/${productId}`)
                    .then(response => {
                        if (!response.ok) {
                            throw new Error(`HTTP error! Status: ${response.status}`);
                        }
                        return response.json();
                    })
                    .then(usedRawMaterialsData => {
                        console.log("Used raw materials data loaded:", usedRawMaterialsData);

                        if (!usedRawMaterialsData.success || !Array.isArray(usedRawMaterialsData.used_raw_materials)) {
                            console.error("Invalid used raw materials data:", usedRawMaterialsData);
                            Swal.fire("Ошибка!", "Не удалось загрузить список использованного сырья.", "error");
                            return;
                        }

                        // Получаем массив ID используемого сырья
                        const usedRawMaterialIds = usedRawMaterialsData.used_raw_materials.map(material => material.id);
                        console.log("Used raw material IDs for product ID:", productId, usedRawMaterialIds);

                        // Фильтруем сырье: исключаем те ингредиенты, которые уже используются
                        const availableRawMaterials = rawMaterialsData.raw_materials.filter(
                            rawMaterial => !usedRawMaterialIds.includes(rawMaterial.id)
                        );
                        console.log("Available raw materials after filtering:", availableRawMaterials);

                        // Создаем выпадающий список для доступного сырья
                        let rawMaterialOptions = availableRawMaterials.map(rawMaterial => {
                            return `<option value="${rawMaterial.id}">${rawMaterial.name}</option>`;
                        }).join("");

                        // Если доступных ингредиентов нет, показываем сообщение
                        if (availableRawMaterials.length === 0) {
                            console.log("No available raw materials for product ID:", productId);
                            Swal.fire({
                                title: "Внимание!",
                                text: "Все доступное сырье уже используется для этого продукта.",
                                icon: "info",
                                timer: 2000,
                                showConfirmButton: false
                            });
                            return;
                        }

                        console.log("Opening modal window for adding an ingredient...");
                        // Открываем модальное окно для добавления ингредиента
                        Swal.fire({
                            title: `Добавление ингредиента для: "${productName}"`,
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
                            <div class="form-group">
    <label for="rawMaterialId" style="display:block; margin-bottom: 5px; font-weight: bold;">Сырье:</label>
    <select id="rawMaterialId" class="input-field">
        ${rawMaterialOptions}
    </select>
</div>

<div class="form-group">
    <label for="quantity" style="display:block; margin-bottom: 5px; font-weight: bold;">Количество:</label>
    <input id="quantity" type="number" class="input-field" placeholder="Введите количество">
</div>
                        `,
                            showCancelButton: true,
                            confirmButtonText: 'Добавить',
                            cancelButtonText: 'Отмена',
                            preConfirm: () => {
                                console.log("Pre-confirm stage: Validating form data...");

                                // Получаем данные из формы
                                const rawMaterialId = document.getElementById('rawMaterialId').value;
                                const quantity = document.getElementById('quantity').value;

                                console.log("Form data submitted:", { rawMaterialId, quantity });

                                // Проверяем, чтобы все поля были заполнены
                                if (!rawMaterialId || !quantity) {
                                    console.warn("Validation failed: Missing required fields.");
                                    Swal.showValidationMessage('Пожалуйста, заполните все поля');
                                    return false;
                                }

                                console.log("Sending data to server for adding ingredient...");
                                // Отправляем данные на сервер
                                return fetch('/ingredients/add', {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json'
                                    },
                                    body: JSON.stringify({
                                        product_id: parseInt(productId),
                                        raw_material_id: parseInt(rawMaterialId),
                                        quantity: parseFloat(quantity)
                                    })
                                })
                                    .then(response => {
                                        if (!response.ok) {
                                            throw new Error(`HTTP error! Status: ${response.status}`);
                                        }
                                        return response.json();
                                    })
                                    .then(data => {
                                        console.log("Server response data for adding ingredient:", data);

                                        if (data.success) {
                                            console.log("Ingredient added successfully.");
                                            return data;
                                        } else {
                                            console.error("Error adding ingredient:", data.error || "Unknown error");
                                            Swal.showValidationMessage(data.error || 'Ошибка при добавлении ингредиента');
                                        }
                                    })
                                    .catch(error => {
                                        console.error("Error during fetch request:", error);
                                        Swal.showValidationMessage('Ошибка при добавлении ингредиента');
                                    });
                            },
                            focusConfirm: false,
                            customClass: {
                                popup: 'popup-class',
                                confirmButton: 'custom-button',
                                cancelButton: 'custom-button'
                            },
                            width: '600px',
                        }).then((result) => {
                            if (result.isConfirmed) {
                                console.log("Ingredient added successfully. Reloading page...");
                                updateIngredientsTable();
                            }
                        });
                    })
                    .catch(error => {
                        console.error("Error fetching used raw materials data:", error);
                        Swal.fire("Ошибка!", "Не удалось загрузить список использованного сырья.", "error");
                    });
            })
            .catch(error => {
                console.error("Error fetching raw materials data:", error);
                Swal.fire("Ошибка!", "Не удалось загрузить список сырья.", "error");
            });
    });
    // Функция для обновления таблицы ингредиентов
    function updateIngredientsTable() {
        const productId = document.getElementById('productSelect').value;

        fetch(`/ingredients/${productId}`)
            .then(response => response.json())
            .then(data => {
                let tableBody = document.getElementById('ingredientsTableBody');
                tableBody.innerHTML = '';

                if (data.ingredients.length === 0) {
                    tableBody.innerHTML = '<div class="table-row"><div class="table-data">Нет данных.</div></div>';
                    return;
                }

                data.ingredients.forEach(ingredient => {
                    let row = `
                        <div class="table-row" id="row-${ingredient.id}">
                            <div class="table-data">${ingredient.material}</div>
                            <div class="table-data">${ingredient.quantity}</div>
                            <div class="table-data action-buttons">
                                <a href="#" class="action-text edit-btn" data-id="${ingredient.id}">
                                    <span>Редактировать</span>
                                    <img src="assets/images/actions/edit.svg" alt="Edit">
                                </a>
                                <a href="#" class="action-text delete-btn" data-id="${ingredient.id}">
                                    <span>Удалить</span>
                                    <img src="assets/images/actions/delete.svg" alt="Delete">
                                </a>
                            </div>
                        </div>
                    `;
                    tableBody.innerHTML += row;
                });
            })
            .catch(error => console.error('Ошибка загрузки данных:', error));
    }

    // Динамическое обновление таблицы при выборе продукта
    document.getElementById('productSelect').addEventListener('change', function () {
        const productId = this.value;
        const productName = this.options[this.selectedIndex].text;

        document.getElementById('selectedProductName').innerText = productName;
        document.getElementById('ingredientSection').style.display = 'block';

        // Обновляем таблицу ингредиентов
        updateIngredientsTable();
    });

    // Стили для SweetAlert2
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
            font-family: 'TildaSans', sans-serif !important;
            font-size: 16px !important;
        }
    `;
    document.head.appendChild(style);
});