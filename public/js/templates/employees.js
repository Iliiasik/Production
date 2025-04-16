// Удаление сотрудника
document.querySelectorAll(".delete-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        const employeeId = this.getAttribute("data-id");

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
                fetch(`/employees/delete/${employeeId}`, { method: "DELETE" })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            // Удаляем строку из DOM
                            const row = document.getElementById(`row-${employeeId}`);
                            if (row) row.remove();

                            // Показываем уведомление об успехе
                            Swal.fire({
                                title: "Удалено!",
                                text: "Сотрудник успешно удален.",
                                icon: "success",
                                timer: 1000,
                                timerProgressBar: true,
                                toast: true,
                                position: "top-end",
                                showConfirmButton: false
                            });

                            // Перезагружаем страницу через 1 секунду
                            setTimeout(() => location.reload(), 1000);
                        } else {
                            Swal.fire("Ошибка!", data.error || "Не удалось удалить сотрудника.", "error");
                        }
                    })
                    .catch(() => {
                        Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
                    });
            }
        });
    });
});
// Редактирование сотрудника
document.querySelectorAll(".edit-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        const employeeId = this.getAttribute("data-id");

        loadData(`/employees/get/${employeeId}`)
            .then(data => {
                if (!data.success || !data.employee) return;

                const emp = data.employee;

                loadData("/positions/list")
                    .then(posData => {
                        if (!posData.success || !posData.positions) return;

                        const posOptions = posData.positions.map(pos => {
                            const selected = pos.id === emp.position_id ? 'selected' : '';
                            return `<option value="${pos.id}" ${selected}>${pos.name}</option>`;
                        }).join("");

                        showModal('Редактировать сотрудника', `
                            <form id="editForm">
                                <table>
                                    <tr>
                                        <th class="table-cell">ФИО</th>
                                        <th class="table-cell">Должность</th>
                                        <th class="table-cell">Оклад</th>
                                        <th class="table-cell">Адрес</th>
                                        <th class="table-cell">Телефон</th>
                                    </tr>
                                    <tr>
                                        <td class="table-cell">
                                            <input type="text" id="fullName" class="input-field" value="${emp.full_name}">
                                        </td>
                                        <td class="table-cell">
                                            <select id="positionId" class="input-field">${posOptions}</select>
                                        </td>
                                        <td class="table-cell">
                                            <input type="number" id="salary" class="input-field" value="${emp.salary}" step="0.01">
                                        </td>
                                        <td class="table-cell">
                                            <input type="text" id="address" class="input-field" value="${emp.address}">
                                        </td>
                                        <td class="table-cell">
                                            <input type="text" id="phone" class="input-field" value="${emp.phone}">
                                        </td>
                                    </tr>
                                </table>
                            </form>
                        `, () => {
                            const full_name = document.getElementById("fullName").value.trim();
                            const position_id = parseInt(document.getElementById("positionId").value);
                            const salary = parseFloat(document.getElementById("salary").value);
                            const address = document.getElementById("address").value.trim();
                            const phone = document.getElementById("phone").value.trim();

                            if (!full_name || isNaN(position_id) || isNaN(salary)) {
                                Swal.showValidationMessage("Заполните все поля");
                                return false;
                            }

                            return fetch(`/employees/edit/${employeeId}`, {
                                method: "POST",
                                headers: { 'Content-Type': 'application/json' },
                                body: JSON.stringify({ full_name, position_id, salary, address, phone })
                            })
                                .then(response => response.json())
                                .then(data => {
                                    if (data.success) {
                                        return data;
                                    } else {
                                        Swal.showValidationMessage(data.error || "Не удалось сохранить изменения");
                                    }
                                })
                                .catch(() => Swal.showValidationMessage("Ошибка при сохранении"));
                        }, '900px').then(result => {
                            if (result.isConfirmed) location.reload();
                        });
                    });
            });
    });
});

// Добавление сотрудника
document.getElementById("addBtn").addEventListener("click", () => {
    fetch("/positions/list")
        .then(response => response.json())
        .then(posData => {
            if (!posData.success || !posData.positions) {
                Swal.fire("Ошибка!", "Не удалось загрузить должности.", "error");
                return;
            }

            const posOptions = posData.positions.map(pos =>
                `<option value="${pos.id}">${pos.name}</option>`
            ).join("");

            showModal('Добавить сотрудника', `
                <form id="addForm">
                    <div class="form-group">
                        <label for="fullName">ФИО:</label>
                        <input type="text" id="fullName" class="input-field" placeholder="Введите ФИО">
                    </div>
                    <div class="form-group">
                        <label for="positionId">Должность:</label>
                        <select id="positionId" class="input-field">${posOptions}</select>
                    </div>
                    <div class="form-group">
                        <label for="salary">Оклад:</label>
                        <input type="number" id="salary" class="input-field" step="0.01" placeholder="Введите оклад">
                    </div>
                    <div class="form-group">
                        <label for="address">Адрес:</label>
                        <input type="text" id="address" class="input-field" placeholder="Введите адрес">
                    </div>
                    <div class="form-group">
                        <label for="phone">Телефон:</label>
                        <input type="text" id="phone" class="input-field" placeholder="Введите телефон">
                    </div>
                </form>
            `, () => {
                const full_name = document.getElementById("fullName").value.trim();
                const position_id = parseInt(document.getElementById("positionId").value);
                const salary = parseFloat(document.getElementById("salary").value);
                const address = document.getElementById("address").value.trim();
                const phone = document.getElementById("phone").value.trim();

                if (!full_name || isNaN(position_id) || isNaN(salary)) {
                    Swal.showValidationMessage("Пожалуйста, заполните все обязательные поля.");
                    return false;
                }

                return fetch("/employees/add", {
                    method: "POST",
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ full_name, position_id, salary, address, phone })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            return data;
                        } else {
                            Swal.showValidationMessage(data.error || "Не удалось добавить сотрудника.");
                        }
                    })
                    .catch(() => Swal.showValidationMessage("Ошибка при добавлении."));
            }, '700px').then(result => {
                if (result.isConfirmed) location.reload();
            });
        });
});