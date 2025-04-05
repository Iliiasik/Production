document.addEventListener("DOMContentLoaded", function () {
    document.getElementById("addBtn").addEventListener("click", function () {
        fetch("/finished-goods/list")
            .then(response => response.json())
            .then(data => {
                if (!data.success) {
                    Swal.fire("Ошибка!", "Не удалось загрузить список продукции", "error");
                    return;
                }
                let finishedGoods = data.finished_goods; // Сохраняем список для быстрого доступа
                let finishedGoodsOptions = finishedGoods.map(
                    fg => `<option value="${fg.id}" data-quantity="${fg.quantity}" data-total="${fg.total_amount}">${fg.name}</option>`
                ).join("");

                fetch("/employees/list")
                    .then(response => response.json())
                    .then(empData => {
                        let employeesOptions = empData.employees.map(
                            emp => `<option value="${emp.id}">${emp.full_name}</option>`
                        ).join("");

                        fetch("/ingredients/list") // Загружаем ингредиенты для отображения необходимого сырья
                            .then(response => response.json())
                            .then(ingredientsData => {
                                if (!ingredientsData.success) {
                                    Swal.fire("Ошибка!", "Не удалось загрузить список ингредиентов", "error");
                                    return;
                                }

                                let ingredients = ingredientsData.ingredients;
                                let today = new Date().toISOString().split("T")[0];

                                Swal.fire({
                                    title: "Производство продукции",
                                    width: "600px",
                                    html: `
                                    <style>
                                        .input-field {
                                            width: 100%;
                                            box-sizing: border-box;
                                            padding: 8px;
                                            border: 1px solid #ccc;
                                            border-radius: 4px;
                                            font-size: 14px;
                                        }
                                        .form-group {
                                            margin-bottom: 16px;
                                        }
                                        .info-box {
                                            background-color: #f2f2f2;
                                            padding: 12px;
                                            margin-bottom: 10px;
                                            border-radius: 8px;
                                            font-size: 18px;
                                            font-weight: bold;
                                        }
                                        .required-materials {
                                            list-style-type: none;
                                            padding: 0;
                                           margin-top: 2%;
                                        }
                                        .required-materials li {
                                            margin-bottom: 8px;
                                            font-size: 16px;
                                            font-weight: 500;
                                            text-align: left;
                                        }
                                        .form-group label {
                                            font-weight: bold;
                                            margin-bottom: 1%;
                                        }
                                    </style>
                                    <div class="form-group">
                                        <label>Продукт:</label>
                                        <select id="productSelect" class="input-field">${finishedGoodsOptions}</select>
                                    </div>
                                    <div class="info-box" id="productInfo">
                                        <div>Текущее количество: <span id="currentQuantity">-</span></div>
                                        <div>Текущая сумма: <span id="currentTotal">-</span></div>
                                    </div>
                                    <div class="info-box" id="rawMaterialInfo">
                                        <div>Необходимое сырье:</div>
                                        <ul id="requiredMaterials" class="required-materials"></ul>
                                    </div>
                                    <div class="form-group">
                                        <label>Количество:</label>
                                        <input type="number" id="quantity" class="input-field" placeholder="Введите количество" min="0">
                                    </div>
                                    <div class="form-group">
                                        <label>Дата производства:</label>
                                        <input type="date" id="productionDate" class="input-field" value="${today}">
                                    </div>
                                    <div class="form-group">
                                        <label>Сотрудник:</label>
                                        <select id="employeeSelect" class="input-field">${employeesOptions}</select>
                                    </div>
                                `,
                                    showCancelButton: true,
                                    confirmButtonText: "Произвести",
                                    cancelButtonText: "Отмена",
                                    focusConfirm: false,
                                    didOpen: () => {
                                        // При открытии окна добавляем обработчик выбора продукта
                                        document.getElementById("productSelect").addEventListener("change", function () {
                                            let selectedOption = this.options[this.selectedIndex];
                                            let productID = parseInt(selectedOption.value, 10);
                                            let quantity = parseFloat(selectedOption.getAttribute("data-quantity")) || 0;
                                            let total = parseFloat(selectedOption.getAttribute("data-total")) || 0;

                                            // Обновляем информацию о продукте
                                            document.getElementById("currentQuantity").innerText = quantity.toFixed(2);
                                            document.getElementById("currentTotal").innerText = total.toFixed(2);

                                            // Отображаем необходимое сырье
                                            let requiredMaterialsList = document.getElementById("requiredMaterials");
                                            requiredMaterialsList.innerHTML = ""; // Очищаем предыдущие данные

                                            // Фильтруем ингредиенты для выбранного продукта
                                            let productIngredients = ingredients.filter(ing => ing.product_id === productID);
                                            if (productIngredients.length === 0) {
                                                let li = document.createElement("li");
                                                li.innerText = "Нет данных о необходимом сырье";
                                                requiredMaterialsList.appendChild(li);
                                                return;
                                            }

                                            productIngredients.forEach(ing => {
                                                let rawMaterial = ing.raw_material; // Сырье уже встроено в объект ингредиента
                                                if (rawMaterial) {
                                                    let li = document.createElement("li");
                                                    li.innerText = `${rawMaterial.name}: необходимо ${ing.quantity.toFixed(2)}, на складе: ${rawMaterial.quantity.toFixed(2)}`;
                                                    requiredMaterialsList.appendChild(li);
                                                } else {
                                                    let li = document.createElement("li");
                                                    li.innerText = `Сырье для ингредиента ID=${ing.raw_material_id} не найдено`;
                                                    requiredMaterialsList.appendChild(li);
                                                }
                                            });
                                        });

                                        // Вызываем событие вручную, чтобы данные сразу загрузились
                                        document.getElementById("productSelect").dispatchEvent(new Event("change"));
                                    },
                                    preConfirm: () => {
                                        const productID = parseInt(document.getElementById("productSelect").value, 10);
                                        const quantity = parseFloat(document.getElementById("quantity").value);
                                        const productionDateInput = document.getElementById("productionDate").value;
                                        const employeeID = parseInt(document.getElementById("employeeSelect").value, 10);

                                        // Проверяем, что все поля заполнены
                                        if (!productID || !quantity || !employeeID || !productionDateInput) {
                                            Swal.showValidationMessage("Все поля должны быть заполнены");
                                            return false;
                                        }

                                        // Проверяем, что количество положительное
                                        if (quantity <= 0) {
                                            Swal.showValidationMessage("Количество должно быть больше нуля");
                                            return false;
                                        }

                                        // Парсим дату без времени
                                        const productionDate = productionDateInput.split("T")[0];

                                        // Проверяем наличие сырья
                                        let productIngredients = ingredients.filter(ing => ing.product_id === productID);
                                        let insufficientMaterials = [];
                                        productIngredients.forEach(ing => {
                                            let rawMaterial = ing.raw_material; // Сырье уже встроено в объект ингредиента
                                            if (!rawMaterial) {
                                                insufficientMaterials.push(`Сырье для ингредиента ID=${ing.raw_material_id} не найдено`);
                                            } else if (rawMaterial.quantity < ing.quantity * quantity) {
                                                insufficientMaterials.push(`${rawMaterial.name} (не хватает ${(ing.quantity * quantity - rawMaterial.quantity).toFixed(2)})`);
                                            }
                                        });

                                        if (insufficientMaterials.length > 0) {
                                            Swal.showValidationMessage(`Недостаточно сырья: ${insufficientMaterials.join(", ")}`);
                                            return false;
                                        }

                                        return { product_id: productID, quantity, production_date: productionDate, employee_id: employeeID };
                                    },
                                    customClass: {
                                        popup: 'popup-class',
                                        confirmButton: 'custom-button',
                                        cancelButton: 'custom-button'
                                    },
                                }).then(result => {
                                    if (result.isConfirmed) {
                                        fetch(`/production/produce/${result.value.product_id}`, {
                                            method: "POST",
                                            headers: { "Content-Type": "application/json" },
                                            body: JSON.stringify(result.value)
                                        })
                                            .then(response => response.json())
                                            .then(data => {
                                                if (data.success) {
                                                    Swal.fire({
                                                        toast: true,
                                                        position: "top-end",
                                                        icon: "success",
                                                        title: "Производство успешно завершено!",
                                                        showConfirmButton: false,
                                                        timer: 1000,
                                                        customClass: "swal-toast"
                                                    });
                                                    setTimeout(() => location.reload(), 1000);
                                                } else {
                                                    Swal.showValidationMessage(data.error);
                                                }
                                            })
                                            .catch(() => Swal.showValidationMessage("Ошибка при отправке данных."));
                                    }
                                });
                            });
                    });
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