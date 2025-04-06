const container = document.body;

container.addEventListener('click', async (e) => {
    const deleteBtn = e.target.closest(".delete-btn");
    const editBtn = e.target.closest(".edit-btn");

    if (deleteBtn) {
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
            customClass: { confirmButton: 'custom-button', cancelButton: 'custom-button' }
        });

        if (confirmed.isConfirmed) {
            try {
                const res = await fetch(`/ingredients/delete/${id}`, { method: "DELETE" });
                const data = await res.json();

                if (data.success) {
                    document.getElementById(`row-${id}`)?.remove();
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
                    updateIngredientsTable();
                } else {
                    Swal.fire("Ошибка!", data.error || "Не удалось удалить запись.", "error");
                }
            } catch {
                Swal.fire("Ошибка!", "Произошла ошибка при удалении.", "error");
            }
        }
    }

    if (editBtn) {
        e.preventDefault();
        const id = editBtn.dataset.id;

        try {
            const res = await fetch(`/ingredients/get/${id}`);
            const { success, ingredient } = await res.json();
            if (!success || !ingredient) throw new Error();

            const [rawRes, usedRes] = await Promise.all([
                fetch("/raw-materials/list").then(r => r.json()),
                fetch(`/ingredients/used-raw-materials/${ingredient.product_id}`).then(r => r.json())
            ]);

            if (!rawRes.success || !usedRes.success) throw new Error();

            const usedIds = usedRes.used_raw_materials.map(m => m.id);
            const options = rawRes.raw_materials
                .filter(m => !usedIds.includes(m.id) || m.id === ingredient.raw_material_id)
                .map(m => `<option value="${m.id}" ${m.id === ingredient.raw_material_id ? 'selected' : ''}>${m.name}</option>`)
                .join("");

            const result = await Swal.fire({
                title: 'Редактировать запись',
                html: `
                    <table>
                        <tr>
                        <th class="table-cell">Сырье</th>
                        <th class="table-cell">Количество</th>
                        </tr>
                        <tr>
                            <td class="table-cell"><select id="rawMaterialId" class="input-field">${options}</select></td>
                            <td class="table-cell"><input id="quantity" type="number" class="input-field" value="${ingredient.quantity}"></td>
                        </tr>
                    </table>`,
                showCancelButton: true,
                confirmButtonText: 'Сохранить изменения',
                cancelButtonText: 'Отмена',
                customClass: {
                    popup: 'popup-class',
                    confirmButton: 'custom-button',
                    cancelButton: 'custom-button'
                },
                preConfirm: async () => {
                    const rawId = document.getElementById('rawMaterialId').value;
                    const qty = document.getElementById('quantity').value;
                    if (!rawId || !qty) return Swal.showValidationMessage('Заполните все поля');

                    const editRes = await fetch(`/ingredients/edit/${id}`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ raw_material_id: +rawId, quantity: +qty })
                    });
                    const editData = await editRes.json();
                    if (!editData.success) Swal.showValidationMessage(editData.error || 'Ошибка');
                    return editData;
                }
            });

            if (result.isConfirmed) updateIngredientsTable();

        } catch {
            Swal.fire("Ошибка!", "Не удалось загрузить данные для редактирования.", "error");
        }
    }
});

document.getElementById('addBtn').addEventListener('click', async () => {
    const productId = document.getElementById('productSelect').value;
    const productName = document.getElementById('productSelect').selectedOptions[0].text;

    try {
        const [rawRes, usedRes] = await Promise.all([
            fetch("/raw-materials/list").then(r => r.json()),
            fetch(`/ingredients/used-raw-materials/${productId}`).then(r => r.json())
        ]);

        if (!rawRes.success || !usedRes.success) throw new Error();

        const usedIds = usedRes.used_raw_materials.map(m => m.id);
        const available = rawRes.raw_materials.filter(m => !usedIds.includes(m.id));

        if (available.length === 0) {
            return Swal.fire({
                title: "Внимание!",
                text: "Все доступное сырье уже используется для этого продукта.",
                icon: "info",
                timer: 2000,
                showConfirmButton: false
            });
        }

        const options = available.map(m => `<option value="${m.id}">${m.name}</option>`).join("");

        const result = await Swal.fire({
            title: `Добавление ингредиента для: "${productName}"`,
            html: `
                <div class="form-group">
                    <label>Сырье:</label>
                    <select id="rawMaterialId" class="input-field">${options}</select>
                </div>
                <div class="form-group">
                    <label>Количество:</label>
                    <input id="quantity" type="number" class="input-field" placeholder="Введите количество">
                </div>`,
            showCancelButton: true,
            confirmButtonText: 'Добавить',
            cancelButtonText: 'Отмена',
            customClass: {
                popup: 'popup-class',
                confirmButton: 'custom-button',
                cancelButton: 'custom-button'
            },
            preConfirm: async () => {
                const rawId = document.getElementById('rawMaterialId').value;
                const qty = document.getElementById('quantity').value;
                if (!rawId || !qty) return Swal.showValidationMessage('Пожалуйста, заполните все поля');

                const res = await fetch('/ingredients/add', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ product_id: +productId, raw_material_id: +rawId, quantity: +qty })
                });

                const data = await res.json();
                if (!data.success) Swal.showValidationMessage(data.error || 'Ошибка при добавлении');
                return data;
            }
        });

        if (result.isConfirmed) updateIngredientsTable();

    } catch (err) {
        Swal.fire("Ошибка!", "Не удалось загрузить данные.", "error");
    }
});

function updateIngredientsTable() {
    const productId = document.getElementById('productSelect').value;
    fetch(`/ingredients/${productId}`)
        .then(res => res.json())
        .then(data => {
            const tbody = document.getElementById('ingredientsTableBody');
            tbody.innerHTML = data.ingredients.length === 0
                ? '<div class="table-row"><div class="table-data">Нет данных.</div></div>'
                : data.ingredients.map(i => `
                    <div class="table-row" id="row-${i.id}">
                        <div class="table-data">${i.material}</div>
                        <div class="table-data">${i.quantity}</div>
                        <div class="table-data action-buttons">
                            <a href="#" class="action-text edit-btn" data-id="${i.id}">
                                <span>Редактировать</span>
                                <img src="assets/images/actions/edit.svg" alt="Edit">
                            </a>
                            <a href="#" class="action-text delete-btn" data-id="${i.id}">
                                <span>Удалить</span>
                                <img src="assets/images/actions/delete.svg" alt="Delete">
                            </a>
                        </div>
                    </div>`).join("");
        });
}

document.getElementById('productSelect').addEventListener('change', function () {
    document.getElementById('selectedProductName').innerText = this.selectedOptions[0].text;
    document.getElementById('ingredientSection').style.display = 'block';
    updateIngredientsTable();
});
