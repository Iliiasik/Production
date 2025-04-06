
    document.body.addEventListener("click", async (e) => {
        const deleteBtn = e.target.closest(".delete-btn");
        if (!deleteBtn) return;

        e.preventDefault();
        const id = deleteBtn.dataset.id;

        const confirmed = await Swal.fire({
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
        });

        if (!confirmed.isConfirmed) return;

        try {
            const res = await fetch(`/purchases/delete/${id}`, { method: "DELETE" });
            const data = await res.json();

            if (data.success) {
                document.getElementById(`row-${id}`)?.remove();

                Swal.fire({
                    title: "Удалено!",
                    text: "Запись успешно удалена.",
                    icon: "success",
                    timer: 1000,
                    toast: true,
                    timerProgressBar: true,
                    position: "top-end",
                    showConfirmButton: false
                });

                setTimeout(() => location.reload(), 1000);
            } else {
                Swal.fire("Ошибка!", data.error || "Не удалось удалить запись.", "error");
            }
        } catch {
            Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
        }
    });
    document.getElementById("addBtn").addEventListener("click", async () => {
        try {
            const [rawData, empData, budgetData] = await Promise.all([
                fetch("/raw-materials/list").then(r => r.json()),
                fetch("/employees/list").then(r => r.json()),
                fetch("/budget/get").then(r => r.json())
            ]);

            if (!rawData.success || !empData.employees) {
                return Swal.fire("Ошибка!", "Не удалось загрузить необходимые данные", "error");
            }

            const rawOptions = rawData.raw_materials.map(rm =>
                `<option value="${rm.id}" data-quantity="${rm.quantity}" data-total="${rm.total_amount}">${rm.name}</option>`
            ).join("");

            const empOptions = empData.employees.map(emp =>
                `<option value="${emp.id}">${emp.full_name}</option>`
            ).join("");

            const today = new Date().toISOString().split("T")[0];
            const budgetAmount = budgetData.success ? budgetData.total_amount : 0;

            const result = await Swal.fire({
                title: "Закупка сырья",
                width: "600px",
                html: `
                <div class="budget-box">Текущий бюджет: ${budgetAmount}</div>
                <div class="form-group">
                    <label>Сырье:</label>
                    <select id="rawMaterialSelect" class="input-field">${rawOptions}</select>
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
                    <select id="employeeSelect" class="input-field">${empOptions}</select>
                </div>
            `,
                showCancelButton: true,
                confirmButtonText: "Закупить",
                cancelButtonText: "Отмена",
                focusConfirm: false,
                customClass: {
                    popup: 'popup-class',
                    confirmButton: 'custom-button',
                    cancelButton: 'custom-button'
                },
                didOpen: () => {
                    const select = document.getElementById("rawMaterialSelect");
                    const qtySpan = document.getElementById("currentQuantity");
                    const totalSpan = document.getElementById("currentTotal");

                    const updateInfo = () => {
                        const opt = select.options[select.selectedIndex];
                        qtySpan.innerText = opt.dataset.quantity || 0;
                        totalSpan.innerText = opt.dataset.total || 0;
                    };

                    select.addEventListener("change", updateInfo);
                    updateInfo();
                },
                preConfirm: () => {
                    const rawId = +document.getElementById("rawMaterialSelect").value;
                    const qty = parseFloat(document.getElementById("quantity").value);
                    const total = parseFloat(document.getElementById("totalAmount").value);
                    const date = document.getElementById("purchaseDate").value;
                    const empId = +document.getElementById("employeeSelect").value;

                    if (!rawId || !qty || !total || !empId || !date) {
                        Swal.showValidationMessage("Все поля должны быть заполнены");
                        return false;
                    }

                    if (total > budgetAmount) {
                        Swal.showValidationMessage("Недостаточно средств в бюджете");
                        return false;
                    }

                    return {
                        raw_material_id: rawId,
                        quantity: qty,
                        total_amount: total,
                        purchase_date: new Date(date).toISOString(),
                        employee_id: empId
                    };
                }
            });

            if (result.isConfirmed) {
                const res = await fetch("/purchases/add", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(result.value)
                });

                const data = await res.json();

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
                    Swal.fire("Ошибка!", data.error || "Не удалось сохранить закупку", "error");
                }
            }
        } catch {
            Swal.fire("Ошибка!", "Произошла ошибка при загрузке данных.", "error");
        }
    });


