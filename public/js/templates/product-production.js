// Производство продукции

document.getElementById("addBtn").addEventListener("click", async function () {
    try {
        const [fgRes, empRes, ingRes] = await Promise.all([
            fetch("/finished-goods/list"),
            fetch("/employees/list"),
            fetch("/ingredients/list"),
        ]);

        const [fgData, empData, ingData] = await Promise.all([
            fgRes.json(),
            empRes.json(),
            ingRes.json(),
        ]);

        if (!fgData.success || !ingData.success) {
            Swal.fire("Ошибка!", "Не удалось загрузить данные", "error");
            return;
        }

        const finishedGoodsOptions = fgData.finished_goods.map(fg =>
            `<option value="${fg.id}" data-quantity="${fg.quantity}" data-total="${fg.total_amount}">${fg.name}</option>`
        ).join("");

        const employeesOptions = empData.employees.map(emp =>
            `<option value="${emp.id}">${emp.full_name}</option>`
        ).join("");

        const ingredients = ingData.ingredients;
        const today = new Date().toISOString().split("T")[0];

        Swal.fire({
            title: "Производство продукции",
            width: "600px",
            html: `
                <div class="form-group">
                    <label>Продукт:</label>
                    <select id="productSelect" class="input-field">${finishedGoodsOptions}</select>
                </div>
                <div class="info-box" id="productInfo">
                    <div>Текущее количество: <span id="currentQuantity">-</span></div>
                    <div>Текущая сумма: <span id="currentTotal">-</span></div>
                </div>
                <div class="info-box">
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
            didOpen: () => updateProductInfo(ingredients),
            preConfirm: () => validateAndPrepareData(ingredients),
            customClass: {
                popup: 'popup-class',
                confirmButton: 'custom-button',
                cancelButton: 'custom-button'
            }
        }).then(result => {
            if (result.isConfirmed) {
                fetch(`/production/produce/${result.value.product_id}`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(result.value)
                })
                    .then(res => res.json())
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

    } catch (err) {
        console.error(err);
        Swal.fire("Ошибка!", "Ошибка при загрузке данных", "error");
    }
});

function updateProductInfo(ingredients) {
    const productSelect = document.getElementById("productSelect");
    const quantitySpan = document.getElementById("currentQuantity");
    const totalSpan = document.getElementById("currentTotal");
    const materialsList = document.getElementById("requiredMaterials");

    const showInfo = () => {
        const selected = productSelect.options[productSelect.selectedIndex];
        const productID = parseInt(selected.value);
        const quantity = parseFloat(selected.getAttribute("data-quantity")) || 0;
        const total = parseFloat(selected.getAttribute("data-total")) || 0;

        quantitySpan.textContent = quantity.toFixed(2);
        totalSpan.textContent = total.toFixed(2);

        materialsList.innerHTML = "";

        const productIngredients = ingredients.filter(i => i.product_id === productID);

        if (productIngredients.length === 0) {
            materialsList.innerHTML = "<li>Нет данных о необходимом сырье</li>";
            return;
        }

        productIngredients.forEach(i => {
            const rm = i.raw_material;
            const line = rm
                ? `${rm.name}: нужно ${i.quantity.toFixed(2)}, на складе: ${rm.quantity.toFixed(2)}`
                : `Сырье для ингредиента ID=${i.raw_material_id} не найдено`;

            const li = document.createElement("li");
            li.textContent = line;
            materialsList.appendChild(li);
        });
    };

    productSelect.addEventListener("change", showInfo);
    showInfo(); // Первичная отрисовка
}

function validateAndPrepareData(ingredients) {
    const productID = parseInt(document.getElementById("productSelect").value, 10);
    const quantity = parseFloat(document.getElementById("quantity").value);
    const productionDate = document.getElementById("productionDate").value;
    const employeeID = parseInt(document.getElementById("employeeSelect").value, 10);

    if (!productID || !quantity || !employeeID || !productionDate) {
        Swal.showValidationMessage("Все поля должны быть заполнены");
        return false;
    }

    if (quantity <= 0) {
        Swal.showValidationMessage("Количество должно быть больше нуля");
        return false;
    }

    const insufficient = ingredients
        .filter(i => i.product_id === productID)
        .filter(i => !i.raw_material || i.raw_material.quantity < i.quantity * quantity)
        .map(i => i.raw_material
            ? `${i.raw_material.name} (не хватает ${(i.quantity * quantity - i.raw_material.quantity).toFixed(2)})`
            : `Сырье для ингредиента ID=${i.raw_material_id} не найдено`
        );

    if (insufficient.length > 0) {
        Swal.showValidationMessage(`Недостаточно сырья: ${insufficient.join(", ")}`);
        return false;
    }

    return {
        product_id: productID,
        quantity,
        production_date: productionDate,
        employee_id: employeeID
    };
}
