document.querySelectorAll(".edit-btn").forEach(button => {
    button.addEventListener("click", function (e) {
        e.preventDefault();
        const id = this.getAttribute("data-id");

        fetch(`/budget/get-row/${id}`)
            .then(r => r.json())
            .then(data => {
                if (!data.success || !data.budget) {
                    Swal.fire("Ошибка!", "Не удалось загрузить данные бюджета.", "error");
                    return;
                }

                const budget = data.budget;

                Swal.fire({
                    title: 'Редактировать бюджет',
                    html: `
                    <form id="editBudgetForm">
                        <table>
                            <tr>
                                <th class="table-cell">Сумма бюджета</th>
                                <th class="table-cell">Наценка (%)</th>
                                <th class="table-cell">Бонус на зарплаты (%)</th>
                            </tr>
                            <tr>
                                <td class="table-cell"><input type="number" id="total_amount" class="input-field" value="${budget.total_amount}" step="0.01"></td>
                                <td class="table-cell"><input type="number" id="markup" class="input-field" value="${budget.markup}" step="0.01"></td>
                                <td class="table-cell"><input type="number" id="salary_bonus" class="input-field" value="${budget.salary_bonus}" step="0.01"></td>
                            </tr>
                        </table>
                    </form>
                `,
                    showCancelButton: true,
                    confirmButtonText: 'Сохранить',
                    cancelButtonText: 'Отмена',
                    preConfirm: () => {
                        const total_amount = parseFloat(document.getElementById("total_amount").value);
                        const markup = parseFloat(document.getElementById("markup").value);
                        const salary_bonus = parseFloat(document.getElementById("salary_bonus").value);

                        if (isNaN(total_amount) || isNaN(markup) || isNaN(salary_bonus)) {
                            Swal.showValidationMessage("Заполните все поля корректно.");
                            return false;
                        }

                        return fetch(`/budget/update/${id}`, {
                            method: "PUT",
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ total_amount, markup, salary_bonus })
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
                    width: '800px'
                }).then(result => {
                    if (result.isConfirmed) location.reload();
                });
            });
    });
});
