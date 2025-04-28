function fetchSalaryData(year, month) {
    fetch(`/salaries/${year}/${month}`)
        .then(res => res.json())
        .then(data => {
            const tbody = document.getElementById('salaryTableBody');
            if (!data || data.length === 0) {
                tbody.innerHTML = '<div class="table-row"><div class="table-data" colspan="8">Нет данных.</div></div>';
                return;
            }

            tbody.innerHTML = data.map(item => {
                const statusIcon = item.is_paid
                    ? `<img src="assets/images/actions/success.svg" alt="Выплачено" class="status-icon" title="Выплачено">`
                    : `<img src="assets/images/actions/error.svg" alt="Не выплачено" class="status-icon" title="Не выплачено">`;

                const editIcon = !item.is_paid
                    ? `<img src="assets/images/actions/edit.svg" class="edit-icon" title="Редактировать" onclick="editSalary(${item.id}, ${item.total_salary})">`
                    : '';

                return `
        <div class="table-row">
            <div class="table-data">${item.employee.full_name}</div>
            <div class="table-data">${item.employee.salary}</div>
            <div class="table-data">${item.purchase_count}</div>
            <div class="table-data">${item.production_count}</div>
            <div class="table-data">${item.sale_count}</div>
            <div class="table-data">${item.total_participation}</div>
            <div class="table-data">${item.bonus}</div>
            <div class="table-data">
                ${item.total_salary}
                ${editIcon}
            </div>
            <div class="table-data">${statusIcon}</div>
        </div>
    `;
            }).join("");
        });
}

function editSalary(recordId, currentSalary) {
    showModal('Редактировать зарплату', `
        <form id="editSalaryForm">
            <table>
                <tr>
                    <th class="table-cell">Сумма зарплаты</th>
                </tr>
                <tr>
                    <td class="table-cell">
                        <input type="number" id="salaryValue" class="input-field" value="${currentSalary}" step="0.01" min="0" placeholder="Введите сумму">
                    </td>
                </tr>
            </table>
        </form>
    `, () => {
        const salary = parseFloat(document.getElementById('salaryValue').value);
        if (isNaN(salary) || salary < 0) {
            Swal.showValidationMessage('Введите корректную сумму зарплаты');
            return false;
        }

        return fetch(`/salaries/edit/${recordId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ total_salary: salary })
        })
            .then(res => {
                if (res.ok) {
                    const year = document.getElementById('yearSelect').value;
                    const month = document.getElementById('monthSelect').value;
                    fetchSalaryData(year, month);
                    updateTotalUnpaidDisplay(year, month);
                } else {
                    Swal.showValidationMessage('Ошибка при обновлении');
                }
            })
            .catch(() => {
                Swal.showValidationMessage('Ошибка при отправке запроса');
            });
    });
}
function updateTotalUnpaidDisplay(year, month) {
    if (!year || !month) {
        document.getElementById('totalUnpaidDisplay').innerText = '';
        return;
    }

    fetch(`/salaries/total-unpaid/${year}/${month}`)
        .then(res => res.json())
        .then(data => {
            if (data.total !== undefined) {
                const totalFormatted = new Intl.NumberFormat('ru-RU', {
                    style: 'currency',
                    currency: 'KGS',
                    minimumFractionDigits: 2
                }).format(data.total);

                document.getElementById('totalUnpaidDisplay').innerText = `К выплате: ${totalFormatted}`;
            } else {
                document.getElementById('totalUnpaidDisplay').innerText = '';
            }
        })
        .catch(() => {
            document.getElementById('totalUnpaidDisplay').innerText = '';
        });
}

// Автоматический расчет и отображение при выборе года/месяца
function calculateAndFetchSalary() {
    const year = document.getElementById('yearSelect').value;
    const month = document.getElementById('monthSelect').value;
    if (!year || !month) return;


    fetch(`/salaries/calculate/${year}/${month}`, { method: 'POST' })
        .then(res => {
            if (res.ok) {
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: 'Зарплата рассчитана успешно',
                    showConfirmButton: false,
                    timer: 2000,
                    timerProgressBar: true
                });
                fetchSalaryData(year, month);
                updateTotalUnpaidDisplay(year, month);
            } else {
                Swal.fire('Ошибка', 'Не удалось рассчитать зарплату', 'error');
            }
        });
}

// Обработка изменения года и месяца
document.getElementById('yearSelect').addEventListener('change', calculateAndFetchSalary);
document.getElementById('monthSelect').addEventListener('change', calculateAndFetchSalary);

// Удалена кнопка showBtn

// Осталась только кнопка выплаты
document.getElementById('payBtn').addEventListener('click', function () {
    const year = document.getElementById('yearSelect').value;
    const month = document.getElementById('monthSelect').value;
    if (!year || !month) return Swal.fire('Ошибка', 'Выберите год и месяц', 'warning');

    fetch(`/salaries/total-unpaid/${year}/${month}`)
        .then(async res => {
            if (!res.ok) {
                Swal.fire('Ошибка', 'Не удалось получить сумму для выплат', 'error');
                return;
            }

            const data = await res.json();
            const total = data.total;

            if (total === 0) {
                Swal.fire({
                    icon: 'info',
                    title: 'Выплата не требуется',
                    text: 'За этот месяц все зарплаты уже выплачены.',
                    confirmButtonText: 'Ок',
                    customClass: {
                        confirmButton: 'custom-button'
                    }
                });
                return;
            }

            Swal.fire({
                title: 'Подтвердите выплату',
                text: `Вы уверены, что хотите выплатить зарплату за выбранный месяц? Сумма: ${total}`,
                icon: 'question',
                showCancelButton: true,
                confirmButtonText: 'Да, выплатить',
                cancelButtonText: 'Отмена',
                customClass: {
                    confirmButton: 'custom-button',
                    cancelButton: 'custom-button'
                }
            }).then(result => {
                if (result.isConfirmed) {
                    paySalaries(year, month);
                }
            });
        })
        .catch(() => {
            Swal.fire('Ошибка', 'Не удалось связаться с сервером', 'error');
        });
});

function paySalaries(year, month) {
    fetch(`/salaries/pay/${year}/${month}`, { method: 'POST' })
        .then(async res => {
            if (res.ok) {
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: 'Зарплаты успешно выплачены',
                    showConfirmButton: false,
                    timer: 2000,
                    timerProgressBar: true
                });
                fetchSalaryData(year, month);
            } else {
                const data = await res.json();
                Swal.fire('Ошибка', data.error || 'Неизвестная ошибка', 'error');
            }
        })
        .catch(() => {
            Swal.fire('Ошибка', 'Не удалось связаться с сервером', 'error');
        });
}
