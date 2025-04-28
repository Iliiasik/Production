document.getElementById('paymentBtn').addEventListener('click', () => {
    const today = new Date().toISOString().split('T')[0];

    // Получаем ID кредита из URL
    const creditId = window.location.pathname.split('/')[2];

    showModal('Добавить платеж', `
        <form id="creditPaymentForm">
            <div class="form-group">
                <label for="paymentDate">Дата платежа:</label>
                <input type="date" id="paymentDate" class="input-field" value="${today}">
            </div>
        </form>
    `, () => {
        const date = document.getElementById('paymentDate').value;

        if (!date) {
            Swal.showValidationMessage('Выберите дату платежа.');
            return false;
        }

        return fetch(`/credits/pay/${creditId}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ date })
        })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    Swal.fire({
                        toast: true,
                        position: 'top-end',
                        icon: 'success',
                        title: 'Платеж выполнен',
                        showConfirmButton: false,
                        timer: 1200,
                        timerProgressBar: true
                    });
                    setTimeout(() => location.reload(), 1000);
                } else {
                    Swal.showValidationMessage(data.error || 'Ошибка при выполнении платежа.');
                }
            })
            .catch(() => Swal.showValidationMessage('Ошибка при отправке запроса.'));
    });
});
