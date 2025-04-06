// Операции удаления, редактирования и добавления записей Finished-goods

document.querySelectorAll(".delete-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        let finishedGoodsId = this.getAttribute("data-id");

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
                fetch(`/finished-goods/delete/${finishedGoodsId}`, {
                    method: "DELETE"
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            const row = document.getElementById(`row-${finishedGoodsId}`);
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
        let finishedGoodsId = this.getAttribute("data-id");

        loadData(`/finished-goods/get/${finishedGoodsId}`)
            .then(data => {
                if (!data || !data.finishedGood) return;

                const finishedGood = data.finishedGood;

                loadData("/units/list")
                    .then(unitsData => {
                        if (!unitsData) return;

                        let unitOptions = unitsData.units.map(unit => {
                            const selected = unit.id === finishedGood.unit_id ? 'selected' : '';
                            return `<option value="${unit.id}" ${selected}>${unit.name}</option>`;
                        }).join("");

                        showModal('Редактировать запись', `
                            <form id="editForm">
                                <table>
                                    <tr>
                                        <th class="table-cell">Название</th>
                                        <th class="table-cell">Единица измерения</th>
                                        <th class="table-cell">Количество</th>
                                        <th class="table-cell">Сумма</th>
                                    </tr>
                                    <tr>
                                        <td class="table-cell">
                                            <input type="text" id="finishedGoodsName" class="input-field" value="${finishedGood.name}">
                                        </td>
                                        <td class="table-cell">
                                            <select id="unitId" class="input-field">
                                                ${unitOptions}
                                            </select>
                                        </td>
                                        <td class="table-cell">
                                            <input type="number" id="quantity" class="input-field" value="${finishedGood.quantity}" step="0.01">
                                        </td>
                                        <td class="table-cell">
                                            <input type="number" id="totalAmount" class="input-field" value="${finishedGood.total_amount}" step="0.01">
                                        </td>
                                    </tr>
                                </table>
                            </form>
                        `, () => {
                            const name = document.getElementById('finishedGoodsName').value.trim();
                            const unitId = document.getElementById('unitId').value;
                            const quantity = parseFloat(document.getElementById('quantity').value);
                            const totalAmount = parseFloat(document.getElementById('totalAmount').value);

                            if (!name || !unitId || isNaN(quantity) || isNaN(totalAmount)) {
                                Swal.showValidationMessage('Заполните все поля');
                                return false;
                            }

                            return fetch(`/finished-goods/edit/${finishedGoodsId}`, {
                                method: 'POST',
                                headers: { 'Content-Type': 'application/json' },
                                body: JSON.stringify({ name, unit_id: parseInt(unitId), quantity, total_amount: totalAmount })
                            })
                                .then(response => response.json())
                                .then(data => {
                                    if (data.success) {
                                        return data;
                                    } else {
                                        Swal.showValidationMessage(data.error || 'Не удалось сохранить изменения');
                                    }
                                })
                                .catch(() => Swal.showValidationMessage('Ошибка при сохранении'));
                        }, '900px').then(result => {
                            if (result.isConfirmed) {
                                location.reload();
                            }
                        });
                    });
            });
    });
});

document.getElementById('addBtn').addEventListener('click', function () {
    fetch("/units/list")
        .then(response => response.json())
        .then(unitsData => {
            if (!unitsData.success) {
                Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
                return;
            }

            let unitOptions = unitsData.units.map(unit => {
                return `<option value="${unit.id}">${unit.name}</option>`;
            }).join("");

            showModal('Добавить запись', `
                <form id="addForm">
                    <div class="form-group">
                        <label for="materialName">Название:</label>
                        <input type="text" id="finishedGoodsName" class="input-field" placeholder="Введите название">
                    </div>

                    <div class="form-group">
                        <label for="unitId">Единица измерения:</label>
                        <select id="unitId" class="input-field">
                            ${unitOptions}
                        </select>
                    </div>
                </form>
            `, () => {
                const name = document.getElementById('finishedGoodsName').value;
                const unitId = document.getElementById('unitId').value;

                if (!name || !unitId) {
                    Swal.showValidationMessage('Пожалуйста, заполните все поля');
                    return false;
                }

                return fetch('/finished-goods/add', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, unit_id: parseInt(unitId) })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            return data;
                        } else {
                            Swal.showValidationMessage(data.error || 'Не удалось добавить запись');
                        }
                    })
                    .catch(() => {
                        Swal.showValidationMessage('Произошла ошибка при добавлении');
                    });
            }).then(result => {
                if (result.isConfirmed) {
                    location.reload(); // Перезагружаем страницу после добавления
                }
            });

        }).catch(error => {
        console.error('Error fetching units data:', error);
        Swal.fire("Ошибка!", "Не удалось загрузить список единиц измерения.", "error");
    });
});