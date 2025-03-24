document.addEventListener("DOMContentLoaded", function () {

    // УДАЛЕНИЕ

    document.querySelectorAll(".delete-btn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            let purchaseId = this.getAttribute("data-id");

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
                    fetch(`/purchases/delete/${purchaseId}`, {
                        method: "DELETE"
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                // Удаляем строку только после успешного удаления на сервере (чтобы не было ошибки окон)
                                const row = document.getElementById(`row-${purchaseId}`);
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
// закупка
    document.getElementById("addBtn").addEventListener("click", function () {
        fetch("/raw-materials/list")
            .then(response => response.json())
            .then(data => {
                if (!data.success) {
                    Swal.fire("Ошибка!", "Не удалось загрузить список сырья", "error");
                    return;
                }

                let rawMaterials = data.raw_materials; // Сохраняем список для быстрого доступа
                let rawMaterialsOptions = rawMaterials.map(
                    rm => `<option value="${rm.id}" data-quantity="${rm.quantity}" data-total="${rm.total_amount}">${rm.name}</option>`
                ).join("");

                fetch("/employees/list")
                    .then(response => response.json())
                    .then(empData => {
                        let employeesOptions = empData.employees.map(
                            emp => `<option value="${emp.id}">${emp.full_name}</option>`
                        ).join("");

                        fetch("/budget/get")
                            .then(response => response.json())
                            .then(budgetData => {
                                let budgetAmount = budgetData.success ? budgetData.total_amount : 0;
                                let today = new Date().toISOString().split("T")[0];

                                Swal.fire({
                                    title: "Закупка сырья",
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
                                    .budget-box, .info-box {
                                        background-color: #f2f2f2;
                                        padding: 12px;
                                        margin-bottom: 10px;
                                        border-radius: 8px;
                                        font-size: 18px;
                                        font-weight: bold;
                                    }
                                    
                                    .form-group label {
                                        font-weight: bold;
                                        margin-bottom: 1%;
                                    }
                                </style>
                                <div class="budget-box">Текущий бюджет: ${budgetAmount}</div>
                                <div class="form-group">
                                    <label>Сырье:</label>
                                    <select id="rawMaterialSelect" class="input-field">${rawMaterialsOptions}</select>
                                </div>
                                <div class="info-box" id="materialInfo">
                                    <div>Текущее количество: <span id="currentQuantity">-</span></div>
                                    <div>Текущая сумма: <span id="currentTotal">-</span></div>
                                </div>
                                <div class="form-group">
                                    <label>Количество:</label>
                                    <input type="number" id="quantity" class="input-field" placeholder="Введите количество">
                                </div>
                                <div class="form-group">
                                    <label>Сумма:</label>
                                    <input type="number" id="totalAmount" class="input-field" placeholder="Введите сумму">
                                </div>
                                <div class="form-group">
                                    <label>Дата закупки:</label>
                                    <input type="date" id="purchaseDate" class="input-field" value="${today}">
                                </div>
                                <div class="form-group">
                                    <label>Сотрудник:</label>
                                    <select id="employeeSelect" class="input-field">${employeesOptions}</select>
                                </div>
                                `,
                                    showCancelButton: true,
                                    confirmButtonText: "Закупить",
                                    cancelButtonText: "Отмена",
                                    focusConfirm: false,
                                    didOpen: () => {
                                        // При открытии окна добавляем обработчик выбора сырья
                                        document.getElementById("rawMaterialSelect").addEventListener("change", function () {
                                            let selectedOption = this.options[this.selectedIndex];
                                            let quantity = selectedOption.getAttribute("data-quantity") || 0;
                                            let total = selectedOption.getAttribute("data-total") || 0;

                                            document.getElementById("currentQuantity").innerText = quantity;
                                            document.getElementById("currentTotal").innerText = total;
                                        });

                                        // Вызываем событие вручную, чтобы данные сразу загрузились
                                        document.getElementById("rawMaterialSelect").dispatchEvent(new Event("change"));
                                    },
                                    preConfirm: () => {
                                        const rawMaterialID = parseInt(document.getElementById("rawMaterialSelect").value, 10);
                                        const quantity = parseFloat(document.getElementById("quantity").value);
                                        const totalAmount = parseFloat(document.getElementById("totalAmount").value);
                                        const purchaseDateInput = document.getElementById("purchaseDate").value;
                                        const employeeID = parseInt(document.getElementById("employeeSelect").value, 10);

                                        const purchaseDate = new Date(purchaseDateInput).toISOString();

                                        if (!rawMaterialID || !quantity || !totalAmount || !employeeID || !purchaseDate) {
                                            Swal.showValidationMessage("Все поля должны быть заполнены");
                                            return false;
                                        }

                                        if (totalAmount > budgetAmount) {
                                            Swal.showValidationMessage("Недостаточно средств в бюджете");
                                            return false;
                                        }

                                        return { raw_material_id: rawMaterialID, quantity, total_amount: totalAmount, purchase_date: purchaseDate, employee_id: employeeID };
                                    },
                                    customClass: {
                                        popup: 'popup-class',
                                        confirmButton: 'custom-button',
                                        cancelButton: 'custom-button'
                                    },
                                }).then(result => {
                                    if (result.isConfirmed) {
                                        fetch("/purchases/add", {
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
                                                        title: "Закупка успешно добавлена!",
                                                        showConfirmButton: false,
                                                        timer: 1000,
                                                        customClass: "swal-toast"
                                                    });

                                                    setTimeout(() => location.reload(), 1000);
                                                } else {
                                                    Swal.fire("Ошибка!", data.error, "error");
                                                }
                                            })
                                            .catch(() => Swal.fire("Ошибка!", "Ошибка при отправке данных.", "error"));
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


