document.querySelectorAll(".delete-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        const positionId = this.getAttribute("data-id");

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
                fetch(`/positions/delete/${positionId}`, { method: "DELETE" })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            // Удаляем строку из DOM
                            const row = document.getElementById(`row-${positionId}`);
                            if (row) row.remove();

                            // Показываем уведомление об успехе
                            Swal.fire({
                                title: "Удалено!",
                                text: "Должность успешно удалена.",
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
                            Swal.fire("Ошибка!", data.error || "Не удалось удалить должность.", "error");
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
        const positionId = this.getAttribute("data-id");

        loadData(`/positions/get/${positionId}`).then(data => {
            if (!data || !data.success) return;

            const position = data.position;
            showModal('Редактировать должность', `
                <form id="editForm">
                    <table>
                        <tr>
                            <th class="table-cell">Название</th>
                        </tr>
                        <tr>
                            <td class="table-cell">
                                <input type="text" id="positionName" class="input-field" value="${position.name}" placeholder="Введите название">
                            </td>
                        </tr>
                    </table>
                </form>
            `, () => {
                const name = document.getElementById('positionName').value.trim();
                if (!name) {
                    Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
                    return false;
                }

                return fetch(`/positions/edit/${positionId}`, {
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
    showModal('Добавить должность', `
        <form id="addForm">
            <div class="form-group">
                <label for="positionName">Название:</label>
                <input type="text" id="positionName" class="input-field" placeholder="Введите название">
            </div>
        </form>
    `, () => {
        const name = document.getElementById('positionName').value.trim();
        if (!name) {
            Swal.showValidationMessage('Пожалуйста, заполните поле: "Название"');
            return false;
        }

        return fetch('/positions/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    location.reload(); // Перезагрузка страницы сразу
                } else {
                    Swal.showValidationMessage(data.error || 'Не удалось добавить должность.');
                }
            })
            .catch(() => Swal.showValidationMessage('Произошла ошибка при добавлении.'));
    });
});