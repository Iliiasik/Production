document.getElementById("addBtn").addEventListener("click", async function () {
    try {
        const [fgRes, empRes, budgetRes] = await Promise.all([
            fetch("/finished-goods/list"),
            fetch("/employees/list"),
            fetch("/markup/get"),
        ]);

        const [fgData, empData, budgetData] = await Promise.all([
            fgRes.json(),
            empRes.json(),
            budgetRes.json(),
        ]);

        if (!fgData.success || !empData.success || !budgetData.success) {
            Swal.fire("Ошибка!", "Не удалось загрузить данные", "error");
            return;
        }

        const finishedGoodsOptions = fgData.finished_goods.map(fg =>
            `<option value="${fg.id}" data-quantity="${fg.quantity}" data-total="${fg.total_amount}">${fg.name}</option>`
        ).join("");

        const employeesOptions = empData.employees.map(emp =>
            `<option value="${emp.id}">${emp.full_name}</option>`
        ).join("");

        const markup = budgetData.markup || 0;
        const today = new Date().toISOString().split("T")[0];

        Swal.fire({
            title: "Продажа продукции",
            width: "650px",
            html: `
                <div class="form-group">
                    <label>Продукт:</label>
                    <select id="productSelect" class="input-field">${finishedGoodsOptions}</select>
                </div>
                <div class="info-box" id="productInfo">
                    <div>Текущее количество: <span id="currentQuantity">-</span></div>
                    <div>Себестоимость всей продукции: <span id="currentTotal">-</span></div>
                    <div>Цена за 1 шт (себестоимость): <span id="costPerUnit">-</span></div>
                    <div>Цена за 1 шт (продажа): <span id="salePerUnit">-</span></div>
                    <div>Стоимость продажи: <span id="saleTotal">-</span></div>
                </div>
                <div class="form-group">
                    <label>Количество:</label>
                    <input type="number" id="quantity" class="input-field" placeholder="Введите количество" min="0">
                </div>
                <div class="form-group">
                    <label>Дата продажи:</label>
                    <input type="date" id="saleDate" class="input-field" value="${today}">
                </div>
                <div class="form-group">
                    <label>Сотрудник:</label>
                    <select id="employeeSelect" class="input-field">${employeesOptions}</select>
                </div>
            `,
            showCancelButton: true,
            confirmButtonText: "Продать",
            cancelButtonText: "Отмена",
            focusConfirm: false,
            didOpen: () => updateProductInfo(markup),
            preConfirm: () => validateAndPrepareData(markup),
            customClass: {
                popup: 'popup-class',
                confirmButton: 'custom-button',
                cancelButton: 'custom-button'
            }
        }).then(result => {
            if (result.isConfirmed) {
                fetch(`/sales/add`, {
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
                                title: "Продажа успешно завершена!",
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

function updateProductInfo(markup) {
    const productSelect = document.getElementById("productSelect");
    const quantityInput = document.getElementById("quantity");

    const quantitySpan = document.getElementById("currentQuantity");
    const totalSpan = document.getElementById("currentTotal");
    const costPerUnitSpan = document.getElementById("costPerUnit");
    const salePerUnitSpan = document.getElementById("salePerUnit");
    const saleTotalSpan = document.getElementById("saleTotal");

    const calculate = () => {
        const selected = productSelect.options[productSelect.selectedIndex];
        const quantityOnStock = parseFloat(selected.getAttribute("data-quantity")) || 0;
        const totalCost = parseFloat(selected.getAttribute("data-total")) || 0;
        const unitCost = quantityOnStock > 0 ? totalCost / quantityOnStock : 0;
        const unitSale = unitCost + unitCost * markup / 100;

        const selectedQty = parseFloat(quantityInput.value) || 0;
        const totalSale = unitSale * selectedQty;

        quantitySpan.textContent = quantityOnStock.toFixed(2);
        totalSpan.textContent = totalCost.toFixed(2);
        costPerUnitSpan.textContent = unitCost.toFixed(2);
        salePerUnitSpan.textContent = unitSale.toFixed(2);
        saleTotalSpan.textContent = totalSale.toFixed(2);
    };

    productSelect.addEventListener("change", calculate);
    quantityInput.addEventListener("input", calculate);
    calculate(); // Первичный расчёт
}

function validateAndPrepareData(markup) {
    const productID = parseInt(document.getElementById("productSelect").value, 10);
    const quantity = parseFloat(document.getElementById("quantity").value);
    const saleDateValue = document.getElementById("saleDate").value;
    const employeeID = parseInt(document.getElementById("employeeSelect").value, 10);

    const selected = document.getElementById("productSelect").selectedOptions[0];
    const quantityOnStock = parseFloat(selected.getAttribute("data-quantity"));
    const totalCost = parseFloat(selected.getAttribute("data-total"));
    const costPerUnit = quantityOnStock > 0 ? totalCost / quantityOnStock : 0;
    const salePerUnit = costPerUnit + costPerUnit * markup / 100;
    const saleTotal = salePerUnit * quantity;

    if (!productID || !quantity || !employeeID || !saleDateValue) {
        Swal.showValidationMessage("Все поля должны быть заполнены");
        return false;
    }

    if (quantity <= 0) {
        Swal.showValidationMessage("Количество должно быть больше нуля");
        return false;
    }

    if (quantity > quantityOnStock) {
        Swal.showValidationMessage("Недостаточно товара на складе");
        return false;
    }

    const saleDate = new Date(saleDateValue).toISOString(); // <-- вот ключ

    return {
        product_id: productID,
        quantity,
        sale_date: saleDate,
        employee_id: employeeID,
        sale_price: saleTotal
    };
}
