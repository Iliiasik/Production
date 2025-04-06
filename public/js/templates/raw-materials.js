// Операции удаления, редактирования и добавления записей Raw-materials

document.querySelectorAll(".delete-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        const id = this.getAttribute("data-id");

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
                fetch(`/raw-materials/delete/${id}`, {
                    method: "DELETE"
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            const row = document.getElementById(`row-${id}`);
                            if (row) row.remove();

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

                            setTimeout(() => location.reload(), 1000);
                        } else {
                            Swal.fire("Ошибка!", data.error || "Не удалось удалить запись.", "error");
                        }
                    })
                    .catch(() => {
                        Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
                    });
            }
        });
    });
});

document.querySelectorAll(".edit-btn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            const id = this.getAttribute("data-id");

            fetch(`/raw-materials/get/${id}`)
                .then(r => r.json())
                .then(data => {
                    if (!data.success || !data.rawmaterial) {
                        Swal.fire("Ошибка!", "Не удалось загрузить данные.", "error");
                        return;
                    }

                    const mat = data.rawmaterial;

                    fetch("/units/list")
                        .then(r => r.json())
                        .then(unitsData => {
                            if (!unitsData.success) {
                                Swal.fire("Ошибка!", "Не удалось загрузить единицы измерения.", "error");
                                return;
                            }

                            const unitOptions = unitsData.units.map(u =>
                                `<option value="${u.id}" ${u.id === mat.unit_id ? 'selected' : ''}>${u.name}</option>`
                            ).join("");

                            Swal.fire({
                                title: 'Редактировать запись',
                                html: `
                                <form id="editForm">
                                    <table>
                                        <tr>
                                            <th class="table-cell">Название</th>
                                            <th class="table-cell">Ед. измерения</th>
                                            <th class="table-cell">Количество</th>
                                            <th class="table-cell">Сумма</th>
                                        </tr>
                                        <tr>
                                            <td class="table-cell"><input type="text" id="name" class="input-field" value="${mat.name}"></td>
                                            <td class="table-cell"><select id="unit" class="input-field">${unitOptions}</select></td>
                                            <td class="table-cell"><input type="number" id="quantity" class="input-field" value="${mat.quantity}" step="0.01"></td>
                                            <td class="table-cell"><input type="number" id="amount" class="input-field" value="${mat.total_amount}" step="0.01"></td>
                                        </tr>
                                    </table>
                                </form>
                            `,
                                showCancelButton: true,
                                confirmButtonText: 'Сохранить',
                                cancelButtonText: 'Отмена',
                                preConfirm: () => {
                                    const name = document.getElementById("name").value.trim();
                                    const unit_id = parseInt(document.getElementById("unit").value);
                                    const quantity = parseFloat(document.getElementById("quantity").value);
                                    const total_amount = parseFloat(document.getElementById("amount").value);

                                    if (!name || isNaN(unit_id) || isNaN(quantity) || isNaN(total_amount)) {
                                        Swal.showValidationMessage("Заполните все поля корректно.");
                                        return false;
                                    }

                                    return fetch(`/raw-materials/edit/${id}`, {
                                        method: "POST",
                                        headers: { 'Content-Type': 'application/json' },
                                        body: JSON.stringify({ name, unit_id, quantity, total_amount })
                                    })
                                        .then(r => r.json())
                                        .then(data => {
                                            if (!data.success) {
                                                Swal.showValidationMessage(data.error || "Не удалось сохранить изменения.");
                                            }
                                            return data;
                                        });
                                },
                                customClass: {
                                    popup: 'popup-class',
                                    confirmButton: 'custom-button',
                                    cancelButton: 'custom-button'
                                },
                                width: '900px'
                            }).then(result => {
                                if (result.isConfirmed) location.reload();
                            });
                        });
                });
        });
    });

document.getElementById("addBtn").addEventListener("click", () => {
        fetch("/units/list")
            .then(r => r.json())
            .then(unitsData => {
                if (!unitsData.success) {
                    Swal.fire("Ошибка!", "Не удалось загрузить единицы измерения.", "error");
                    return;
                }

                const unitOptions = unitsData.units.map(u =>
                    `<option value="${u.id}">${u.name}</option>`
                ).join("");

                Swal.fire({
                    title: "Добавить запись",
                    html: `
                    <form id="addForm">
                        <div class="form-group">
                            <label for="name">Название:</label>
                            <input type="text" id="name" class="input-field" placeholder="Введите название">
                        </div>
                        <div class="form-group">
                            <label for="unit">Единица измерения:</label>
                            <select id="unit" class="input-field">${unitOptions}</select>
                        </div>
                    </form>
                `,
                    showCancelButton: true,
                    confirmButtonText: "Добавить",
                    cancelButtonText: "Отмена",
                    preConfirm: () => {
                        const name = document.getElementById("name").value.trim();
                        const unit_id = parseInt(document.getElementById("unit").value);

                        if (!name || isNaN(unit_id)) {
                            Swal.showValidationMessage("Заполните все поля.");
                            return false;
                        }

                        return fetch("/raw-materials/add", {
                            method: "POST",
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ name, unit_id })
                        })
                            .then(r => r.json())
                            .then(data => {
                                if (!data.success) {
                                    Swal.showValidationMessage(data.error || "Не удалось добавить запись.");
                                }
                                return data;
                            });
                    },
                    customClass: {
                        popup: 'popup-class',
                        confirmButton: 'custom-button',
                        cancelButton: 'custom-button'
                    },
                    width: '600px'
                }).then(result => {
                    if (result.isConfirmed) location.reload();
                });
            });
    });


