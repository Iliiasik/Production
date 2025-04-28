document.getElementById('addBtn').addEventListener('click', function () {
    // Получаем сегодня в формате YYYY-MM-DD
    const today = new Date().toISOString().split('T')[0];

    showModal('Добавить запись', `
        <form id="addCreditForm">
            <div class="form-group">
                <label for="amount">Сумма кредита:</label>
                <input type="number" id="amount" class="input-field" placeholder="Введите сумму">
            </div>
            <div class="form-group">
                <label for="startDate">Дата получения:</label>
                <input type="date" id="startDate" class="input-field" value="${today}">
            </div>
            <div class="form-group">
                <label for="termYears">Срок (в годах):</label>
                <input type="number" id="termYears" class="input-field" placeholder="Введите срок">
            </div>
            <div class="form-group">
                <label for="annualRate">Годовой процент (%):</label>
                <input type="number" id="annualRate" step="0.01" class="input-field" placeholder="Введите процент">
            </div>
            <div class="form-group">
                <label for="penaltyRate">Пеня (% в день):</label>
                <input type="number" id="penaltyRate" step="0.01" class="input-field" placeholder="Введите пеню">
            </div>
        </form>
    `, () => {
        const amount      = parseFloat(document.getElementById('amount').value.trim());
        const startDate   = document.getElementById('startDate').value;
        const termYears   = parseInt(document.getElementById('termYears').value.trim());
        const annualRate  = parseFloat(document.getElementById('annualRate').value.trim());
        const penaltyRate = parseFloat(document.getElementById('penaltyRate').value.trim());

        if (isNaN(amount) || !startDate || isNaN(termYears) || isNaN(annualRate) || isNaN(penaltyRate)) {
            Swal.showValidationMessage('Пожалуйста, заполните все поля корректно.');
            return false;
        }

        return fetch('/credits/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                amount,
                start_date: startDate,
                term_years: termYears,
                annual_rate: annualRate,
                penalty_rate: penaltyRate
            })
        })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    Swal.fire({
                        toast: true,
                        position: 'top-end',
                        icon: 'success',
                        title: 'Кредит добавлен',
                        showConfirmButton: false,
                        timerProgressBar: true,
                        timer: 1000
                    });
                    setTimeout(() => location.reload(), 1000);
                } else {
                    Swal.showValidationMessage(data.error || 'Не удалось добавить кредит.');
                }
            })
            .catch(() => Swal.showValidationMessage('Произошла ошибка при добавлении.'));
    });
});
