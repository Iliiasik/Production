// Карта маршрутов для разделов
const reportRoutes = {
    sales: '/reports/sales',
    productions: '/reports/productions',
    purchases: '/reports/purchases',
    salaries: '/reports/salaries',
    payments: '/reports/payments'
};

// Главная кнопка генерации
const generateBtn = document.getElementById('generate-report-btn');

generateBtn.addEventListener('click', async () => {
    const section = document.getElementById('report-select').value;
    const startDate = document.getElementById('start-date').value;
    const endDate = document.getElementById('end-date').value;

    // Валидация
    if (!section || !startDate || !endDate) {
        Swal.fire({
            icon: 'warning',
            title: 'Заполните все поля!',
            text: 'Пожалуйста, выберите раздел и диапазон дат.',
            customClass: {
                confirmButton: 'custom-button'
            }
        });
        return;
    }

    const url = reportRoutes[section];
    if (!url) {
        Swal.fire({
            icon: 'error',
            title: 'Неизвестный раздел',
            text: 'Не удалось определить маршрут для отчёта.',
            customClass: {
                confirmButton: 'custom-button'
            }
        });
        return;
    }

    Swal.fire({
        title: 'Генерация отчёта...',
        allowOutsideClick: false,
        didOpen: () => Swal.showLoading()
    });

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ start_date: startDate, end_date: endDate })
        });

        if (!response.ok) throw new Error(`Ошибка: ${response.status}`);

        const html = await response.text();
        insertReportTable(html);

        Swal.close();
    } catch (err) {
        console.error('Ошибка генерации отчёта:', err);
        Swal.fire({
            icon: 'error',
            title: 'Ошибка',
            text: 'Не удалось загрузить отчёт.',
            customClass: {
                confirmButton: 'custom-button'
            }
        });
    }
});

function insertReportTable(html) {
    const container = document.getElementById('report-result');
    container.innerHTML = html;

    const exportBtn = document.getElementById('export-btn');
    if (exportBtn) {
        // Показываем контейнер с кнопкой экспорта
        exportBtn.parentElement.style.display = 'inline-block';

        // Назначаем обработчик для модального окна выбора формата
        exportBtn.addEventListener('click', () => {
            showExportModal();
        });
    }
}

// Модифицированная функция showExportModal
function showExportModal() {
    Swal.fire({
        title: 'Выберите формат экспорта',
        html: `
            <div style="display: flex; justify-content: center; gap: 1rem; margin-top: 1rem;">
                <button class="profile-button" onclick="exportReport('docx')">
                    <img src="/assets/images/logos/word.svg" alt="Word" width="32"><br>Word
                </button>
                <button class="profile-button" onclick="exportReport('xlsx')">
                    <img src="/assets/images/logos/excel.svg" alt="Excel" width="32"><br>Excel
                </button>
                <button class="profile-button" onclick="exportReport('pdf')">
                    <img src="/assets/images/logos/pdf.svg" alt="PDF" width="32"><br>PDF
                </button>
            </div>
        `,
        showConfirmButton: false
    });
}

// Универсальная функция экспорта
async function exportReport(format) {
    Swal.fire({
        title: 'Подготовка отчёта...',
        allowOutsideClick: false,
        didOpen: () => Swal.showLoading()
    });

    try {
        const reportContent = document.getElementById('report-result').innerHTML;
        const reportTitle = document.querySelector('.report-title')?.textContent || 'Отчёт';
        const period = document.querySelector('.report-period')?.textContent || '';

        const response = await fetch('/reports/export', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                content: reportContent,
                title: reportTitle,
                period: period,
                format: format
            })
        });

        if (!response.ok) throw new Error(`Ошибка: ${response.status}`);

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${reportTitle.replace(/ /g, '_')}_${new Date().toISOString().slice(0,10)}.${format}`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);

        Swal.close();
    } catch (err) {
        console.error('Ошибка экспорта отчёта:', err);
        Swal.fire({
            icon: 'error',
            title: 'Ошибка',
            text: 'Не удалось экспортировать отчёт.',
            customClass: {
                confirmButton: 'custom-button'
            }
        });
    }
}