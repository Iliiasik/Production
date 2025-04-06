// Операции удаления, редактирования и добавления записей Units

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
                            const row = document.getElementById(`row-${unitId}`);
                            if (row) row.remove(); // Удаление строки
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
                            setTimeout(() => location.reload(), 1000); // Перезагрузка страницы через 1 секунду
                        } else {
                            Swal.fire("Ошибка!", data.error || "Не удалось удалить запись.", "error");
                        }
                    })
                    .catch(error => {
                        Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
                    });
            }
        });
    });
});

document.querySelectorAll(".edit-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        let unitId = this.getAttribute("data-id");

        loadData(`/units/get/${unitId}`).then(data => {
            if (!data || !data.success) return;

            const unit = data.unit;
            showModal('Редактировать запись', `
                <form id="editForm">
                    <table>
                        <tr>
                            <th class="table-cell">Название</th>
                        </tr>
                        <tr>
                            <td class="table-cell">
                                <input type="text" id="unitName" class="input-field" value="${unit.name}" placeholder="Введите название">
                            </td>
                        </tr>
                    </table>
                </form>
            `, () => {
                const name = document.getElementById('unitName').value.trim();
                if (!name) {
                    Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
                    return false;
                }

                return fetch(`/units/edit/${unitId}`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            location.reload(); // Перезагрузка страницы сразу
                        } else {
                            Swal.showValidationMessage(data.error || 'Не удалось сохранить изменения');
                        }
                    })
                    .catch(() => Swal.showValidationMessage('Ошибка при сохранении'));
            });
        });
    });
});

document.getElementById('addBtn').addEventListener('click', function () {
    loadData("/units/list").then(unitsData => {
        if (!unitsData || !unitsData.success) {
            Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
            return;
        }

        let unitOptions = unitsData.units.map(unit => {
            return `<option value="${unit.id}">${unit.name}</option>`;
        }).join("");

        showModal('Добавить запись', `
            <form id="addForm">
                <div>
                    <label for="unitName">Название:</label>
                    <input type="text" id="unitName" class="input-field" placeholder="Введите название">
                </div>
            </form>
        `, () => {
            const name = document.getElementById('unitName').value.trim();
            if (!name) {
                Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
                return false;
            }

            return fetch('/units/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        location.reload(); // Перезагрузка страницы сразу
                    } else {
                        Swal.showValidationMessage(data.error || 'Не удалось добавить запись');
                    }
                })
                .catch(() => Swal.showValidationMessage('Произошла ошибка при добавлении'));
        });
    });
});
