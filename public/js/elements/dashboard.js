const main = document.getElementById('user-dashboard');
const fullName = main.dataset.fullname;
const position = main.dataset.position;
const salary = main.dataset.salary;
const username = main.dataset.username;
const ispasswordchanged = main.dataset.ispasswordchanged;

// Блок предупреждения пароля
const passwordChangeWarning = ispasswordchanged === "false" ? `
  <div class="password-warning">
    <img src="assets/images/actions/warning.svg" alt="Warning" class="warning-icon">
    <p><strong>Смените пароль для безопасности.</strong></p>
  </div>
` : '';

// Блок разметки дешборда
const html = `
<section class="cards-wrapper">
  <div class="card card--large">
    <div class="profile-info">
      <h2>Добро пожаловать, ${fullName}!</h2>
      <p class="username-text">@${username}</p>

      <div class="glass-cards-container">
        <div class="glass-card">
          <h3>Должность</h3>
          <p class="position-text">${position}</p>
        </div>

        <div class="glass-card">
          <h3>Зарплата</h3>
          <p class="salary-text">${salary}</p>
        </div>
      </div>

    </div>
  </div>

  <div class="card-col">
    <div class="card card--small">
      <div class="actions-list">
        <button id="change-password-button" class="profile-button">
          <div class="sign"></div>
          <div class="text">Сменить пароль</div>
        </button>
        ${passwordChangeWarning}  <!-- Здесь предупреждение -->
        
        <button id="logout-button" class="profile-button">
          <div class="sign"></div>
          <div class="text">Выйти</div>
        </button>
      </div>
    </div>

    <div class="card card--small">
      <div class="actions-list">
        <button id="show-permissions-button" class="profile-button">Ваши доступы</button>
        ${position === 'Админ' ? `
          <button id="manage-access-button" class="profile-button">Управление доступами</button>
        ` : ''}
        ${position !== 'Директор' ? `
  <button id="taskboard-button" class="profile-button">Мои задачи</button>
` : ''}
${position === 'Директор' ? `
      <button id="manage-tasks-button" class="profile-button">Управление задачами</button>
` : ''}

      </div>
    </div>
  </div>
  
</section>
`;
main.innerHTML = html;

// Блок учета и контроля задач
document.getElementById('taskboard-button')?.addEventListener('click', async () => {
    try {
        const response = await fetch('/tasks/my');
        const data = await response.json();
        const tasks = data.tasks;

        let html = `<h4 style="margin-bottom: 1rem;">Мои задачи</h4>`;
        if (!tasks || tasks.length === 0) {
            html += `<p>У вас нет активных задач</p>`;
        } else {
            html += `<div class="task-cards-container">`;

            tasks.forEach(task => {
                const dueDate = task.due_date
                    ? new Date(task.due_date).toLocaleDateString()
                    : '—';

                const statusColor = getStatusColor(task.status);

                html += `
                    <div class="task-card">
                        <div class="task-card-header">
                            <div class="status-dot" style="background-color: ${statusColor};"></div>
                            <strong>${task.title || 'Без названия'}</strong>
                        </div>
                        <div class="task-card-body">
                            <p>${task.description || 'Нет описания'}</p>
                            <div class="task-meta">
                                <span><strong>Статус:</strong> ${task.status || 'Неизвестен'}</span>
                                <span><strong>Срок:</strong> ${dueDate}</span>
                            </div>
                        </div>
                    </div>
                `;
            });

            html += `</div>`;
        }

        Swal.fire({
            title: 'Ваш таскборд',
            html: html,
            width: 800,
            confirmButtonText: 'Закрыть',
            customClass: {
                confirmButton: 'custom-button'
            }
        });

    } catch (error) {
        console.error('Ошибка загрузки задач:', error);
        Swal.fire({
            title: 'Ошибка',
            text: 'Не удалось загрузить задачи',
            icon: 'error',
            customClass: {
                confirmButton: 'custom-button'
            }
        });
    }
    function getStatusColor(status) {
        switch (status) {
            case 'Новая':
                return '#007bff';
            case 'В работе':
                return '#ffc107';
            case 'Завершена':
                return '#28a745';
            default:
                return '#6c757d';
        }
    }

});
document.getElementById('manage-tasks-button')?.addEventListener('click', async () => {
    try {
        const response = await fetch('/employees/list');
        const data = await response.json();

        if (!data.success) {
            throw new Error('Не удалось загрузить сотрудников');
        }

        const employees = data.employees;

        // Удаляем старую панель, если она есть
        const existing = document.getElementById('permissions-panel');
        if (existing) existing.remove();

        // Создаем панель
        const panel = document.createElement('div');
        panel.id = 'permissions-panel';
        panel.className = 'permissions-panel'; // Применяем тот же класс, что и для панели доступа
        panel.innerHTML = `
            <h4>Задачи</h4>
                        <div class="select-mode">
            <button id="add-task-btn" class="profile-button">Добавить задачу</button>
            <button id="view-tasks-btn" class="profile-button">Просмотр задач</button>
            </div>
        `;
        document.getElementById('user-dashboard').appendChild(panel);

        // Обработчик для добавления задачи через SweetAlert
        document.getElementById('add-task-btn').addEventListener('click', () => {
            fetch("/employees/list")
                .then(res => res.json())
                .then(empData => {
                    if (!empData.success || !empData.employees) {
                        Swal.fire("Ошибка!", "Не удалось загрузить сотрудников.", "error");
                        return;
                    }

                    const employeeOptions = empData.employees.map(emp =>
                        `<option value="${emp.id}">${emp.full_name}</option>`
                    ).join("");

                    showModal('Добавить запись', `
                <div class="form-group">
                    <label for="assign-to">Сотрудник:</label>
                    <select id="assign-to" class="input-field">${employeeOptions}</select>
                </div>
                <div class="form-group">
                    <label for="task-title">Заголовок:</label>
                    <input type="text" id="task-title" class="input-field" placeholder="Введите заголовок">
                </div>
                <div class="form-group">
                    <label for="task-desc">Описание:</label>
                    <textarea id="task-desc" class="input-field" placeholder="Введите описание"></textarea>
                </div>
                <div class="form-group">
                    <label for="task-date">Срок:</label>
                    <input type="date" id="task-date" class="input-field">
                </div>
            `, () => {
                        const title = document.getElementById("task-title").value.trim();
                        const description = document.getElementById("task-desc").value.trim();
                        const assignedTo = parseInt(document.getElementById("assign-to").value);
                        const dueDateRaw = document.getElementById("task-date").value;

                        if (!title || isNaN(assignedTo) || !dueDateRaw) {
                            Swal.showValidationMessage("Пожалуйста, заполните все обязательные поля корректно.");
                            return false;
                        }

                        return fetch("/tasks", {
                            method: "POST",
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                title,
                                description,
                                assigned_to: assignedTo,
                                due_date: dueDateRaw
                            })
                        })
                            .then(response => response.json())
                            .then(data => {
                                if (data.success) {
                                    // Показываем тост в правом верхнем углу
                                    Swal.fire({
                                        position: 'top-end',
                                        icon: 'success',
                                        title: 'Задача успешно добавлена',
                                        showConfirmButton: false,
                                        timer: 3000,
                                        toast: true,
                                        background: '#f8f9fa',
                                        timerProgressBar: true
                                    });
                                    return true;
                                } else {
                                    throw new Error(data.error || "Ошибка сервера");
                                }
                            })
                            .catch(error => {
                                Swal.showValidationMessage(error.message);
                                return false;
                            });
                    }, '700px');
                })
                .catch(error => {
                    Swal.fire("Ошибка!", "Не удалось загрузить данные: " + error.message, "error");
                });
        });
        document.getElementById('view-tasks-btn').addEventListener('click', async () => {
            try {
                console.log("Запрос задач...");
                const response = await fetch('/tasks/list');

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();
                console.log("Полученные данные:", data);

                if (!data || !data.success || !Array.isArray(data.tasks)) {
                    throw new Error("Некорректная структура данных");
                }

                const tasks = data.tasks;
                if (tasks.length === 0) {
                    alert("Нет задач для отображения");
                    return;
                }

                // Удаляем старую панель
                const existingPanel = document.getElementById('tasks-panel');
                if (existingPanel) existingPanel.remove();

                // Создаем новую панель
                const panel = document.createElement('div');
                panel.id = 'permissions-panel';
                panel.className = 'permissions-panel';

                let tasksHtml = `<h4>Список задач (${tasks.length}):</h4>`;

                // Группировка по сотрудникам
                const groupedTasks = tasks.reduce((acc, task) => {
                    const key = task.employee?.full_name || 'Не назначено';
                    if (!acc[key]) {
                        acc[key] = {
                            employee: key,
                            tasks: []
                        };
                    }
                    acc[key].tasks.push(task);
                    return acc;
                }, {});

                // Генерация HTML
                for (const group of Object.values(groupedTasks)) {
                    tasksHtml += `
            <div class="permission-category-card">
                <div class="permission-category-header">
                    <span>${group.employee}</span>
                    <svg class="arrow-icon" viewBox="0 0 24 24" width="18" height="18">
                        <path fill="currentColor" d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6z"/>
                    </svg>
                </div>
                <ul class="permission-list hidden">
        `;

                    group.tasks.forEach(task => {
                        const dueDate = task.due_date
                            ? new Date(task.due_date).toLocaleDateString()
                            : 'не указана';

                        tasksHtml += `
                <li class="permission-item">
                    <span class="checkmark">
                        <svg viewBox="0 0 24 24" width="16" height="16">
                            <circle cx="12" cy="12" r="10" fill="${getStatusColor(task.status)}" />
                        </svg>
                    </span>
                    <div class="task-info">
                        <span class="permission-desc">${task.title || 'Без названия'}</span><br>
                        <div class="task-details">
                            <small><strong>Описание:</strong> ${task.description || 'Не указано'}</small><br>
                            <small><strong>Срок:</strong> ${dueDate}</small><br>
                            <small><strong>Статус:</strong> ${task.status || 'Неизвестен'}</small>
                        </div>
                    </div>
                </li>
            `;
                    });

                    tasksHtml += `</ul></div>`;
                }

                tasksHtml += `
        <div class="actions-center">
            <button id="close-tasks" class="profile-button">Закрыть</button>
        </div>
        `;

                panel.innerHTML = tasksHtml;
                document.getElementById('user-dashboard').appendChild(panel);

                // Обработчики событий
                document.getElementById('close-tasks').addEventListener('click', () => panel.remove());

                panel.querySelectorAll('.permission-category-header').forEach(header => {
                    header.addEventListener('click', () => {
                        const list = header.nextElementSibling;
                        list.classList.toggle('hidden');
                        header.querySelector('.arrow-icon').classList.toggle('rotated');
                    });
                });

            } catch (error) {
                console.error('Ошибка:', error);
                alert(`Ошибка: ${error.message}`);
            }
        });
        function getStatusColor(status) {
            switch (status) {
                case 'Новая':
                    return '#007bff';
                case 'В работе':
                    return '#ffc107';
                case 'Завершена':
                    return '#28a745';
                default:
                    return '#6c757d';
            }
        }


    } catch (error) {
        console.error('Ошибка загрузки сотрудников:', error);
        Swal.fire('Ошибка', 'Не удалось загрузить сотрудников', 'error');
    }
});

// Блок отображения доступов пользователя
const showPermissionsButton = document.getElementById('show-permissions-button');
showPermissionsButton.addEventListener('click', async () => {
    try {
        const response = await fetch('/user/permissions');
        if (!response.ok) {
            alert('Не удалось загрузить разрешения');
            return;
        }

        const data = await response.json();
        const permissions = data.permissions;

        let existingPanel = document.getElementById('permissions-panel');
        if (existingPanel) {
            existingPanel.remove();
        }

        const panel = document.createElement('div');
        panel.id = 'permissions-panel';
        panel.className = 'permissions-panel';

        let permissionsHtml = `<h4>Ваши доступы:</h4>`;

        const groupedPermissions = permissions.reduce((acc, permission) => {
            if (!acc[permission.category]) {
                acc[permission.category] = [];
            }
            acc[permission.category].push(permission);
            return acc;
        }, {});

        for (const category in groupedPermissions) {
            permissionsHtml += `
                <div class="permission-category-card">
                    <div class="permission-category-header">
                        <span>${category}</span>
                        <svg class="arrow-icon" viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6z"/></svg>
                    </div>
                    <ul class="permission-list hidden">
            `;

            groupedPermissions[category].forEach(permission => {
                permissionsHtml += `
                    <li class="permission-item">
                        <span class="checkmark">
                            <svg viewBox="0 0 24 24" width="16" height="16">
                                <path fill="#2ecc71" d="M20.285,6.709c-0.391-0.391-1.023-0.391-1.414,0L9,16.58l-3.871-3.871c-0.391-0.391-1.023-0.391-1.414,0s-0.391,1.023,0,1.414l4.578,4.578c0.391,0.391,1.023,0.391,1.414,0l10.578-10.578C20.676,7.732,20.676,7.1,20.285,6.709z"/>
                            </svg>
                        </span>
                        <span class="permission-desc">${permission.description}</span>
                    </li>
                `;
            });

            permissionsHtml += `</ul></div>`;
        }

        permissionsHtml += `
            <div class="actions-center">
                <button id="close-permissions" class="profile-button">Закрыть</button>
            </div>
        `;

        panel.innerHTML = permissionsHtml;

        const main = document.getElementById('user-dashboard');
        main.appendChild(panel);

        document.getElementById('close-permissions').addEventListener('click', () => {
            panel.remove();
        });

        // Добавляем логику сворачивания/разворачивания
        const headers = panel.querySelectorAll('.permission-category-header');
        headers.forEach(header => {
            header.addEventListener('click', () => {
                const list = header.nextElementSibling;
                list.classList.toggle('hidden');
                header.querySelector('.arrow-icon').classList.toggle('rotated');
            });
        });

    } catch (error) {
        console.error('Ошибка загрузки разрешений:', error);
        alert('Ошибка загрузки разрешений');
    }
});

// Блок управления доступами
const manageAccessButton = document.getElementById('manage-access-button');
if (manageAccessButton) {
    manageAccessButton.addEventListener('click', () => {
        // Удаляем старую панель
        const existing = document.getElementById('permissions-panel');
        if (existing) existing.remove();

        // Рисуем контейнер
        const panel = document.createElement('div');
        panel.id = 'permissions-panel';
        panel.className = 'permissions-panel';
        panel.innerHTML = `
            <h4>Управление доступами</h4>
            <div class="select-mode">
                <button id="roles-mode-button" class="profile-button">Должности</button>
                <button id="users-mode-button" class="profile-button">Пользователи</button>
            </div>
            <div id="management-content" class="management-content hidden"></div>
            <div class="actions-center">
                <button id="close-permissions" class="profile-button">Закрыть</button>
                <button id="save-role-permissions-button" class="profile-button hidden">Сохранить должность</button>
                <button id="save-user-permissions-button" class="profile-button hidden">Сохранить пользователя</button>
            </div>
        `;
        document.getElementById('user-dashboard').appendChild(panel);

        // Обработчик закрытия панели
        document.getElementById('close-permissions')
            .addEventListener('click', () => panel.remove());

        const managementContent = document.getElementById('management-content');

        // --- Обработчики для кнопок сохранения ---
        document.getElementById('save-role-permissions-button').addEventListener('click', async () => {
            const roleId = document.getElementById('role-select')?.value;
            if (!roleId) {
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'error',
                    title: 'Выберите должность',
                    showConfirmButton: false,
                    timer: 3000
                });
                return;
            }

            const ids = Array.from(
                document.querySelectorAll('#role-permissions input:checked')
            ).map(ch => Number(ch.dataset.id));

            try {
                const res = await fetch(`/admin/roles/${roleId}/permissions/update`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ permission_ids: ids })
                });

                if (!res.ok) throw new Error(await res.text());

                const data = await res.json();
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: 'Должность обновлена',
                    showConfirmButton: false,
                    timer: 3000
                });
            } catch (error) {
                console.error('Ошибка сохранения должности:', error);
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'error',
                    title: 'Ошибка сохранения',
                    text: error.message || 'Неизвестная ошибка',
                    showConfirmButton: false,
                    timer: 5000
                });
            }
        });

        document.getElementById('save-user-permissions-button').addEventListener('click', async () => {
            const userId = document.getElementById('user-select')?.value;
            if (!userId) {
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'error',
                    title: 'Выберите пользователя',
                    showConfirmButton: false,
                    timer: 3000
                });
                return;
            }

            const ids = Array.from(
                document.querySelectorAll('#user-permissions input:checked:not(:disabled)')
            ).map(ch => Number(ch.dataset.id));

            try {
                const res = await fetch(`/admin/users/${userId}/permissions/update`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ permission_ids: ids })
                });

                if (!res.ok) throw new Error(await res.text());

                const data = await res.json();
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: 'Пользователь обновлён',
                    showConfirmButton: false,
                    timer: 3000
                });
            } catch (error) {
                console.error('Ошибка сохранения пользователя:', error);
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'error',
                    title: 'Ошибка сохранения',
                    text: error.message || 'Неизвестная ошибка',
                    showConfirmButton: false,
                    timer: 5000
                });
            }
        });

        // --- РЕЖИМ РОЛЕЙ ---
        document.getElementById('roles-mode-button').addEventListener('click', async () => {
            document.getElementById('save-role-permissions-button').classList.remove('hidden');
            document.getElementById('save-user-permissions-button').classList.add('hidden');

            managementContent.classList.remove('hidden');
            managementContent.innerHTML = '<p>Загрузка ролей...</p>';

            try {
                const respRoles = await fetch('/admin/roles');
                if (!respRoles.ok) throw new Error('Ошибка загрузки должностей');

                const { positions } = await respRoles.json();
                const opts = positions.map(r => `<option value="${r.id}">${r.name}</option>`).join('');
                managementContent.innerHTML = `
                    <select id="role-select" class="glass-select">
                        <option value="" disabled selected>Выберите должность</option>
                        ${opts}
                    </select>
                    <div id="role-permissions" class="permissions-list"></div>
                `;

                const roleSelect = document.getElementById('role-select');
                roleSelect.addEventListener('change', async () => {
                    const roleId = roleSelect.value;
                    if (!roleId) return;

                    try {
                        const [allR, roleR] = await Promise.all([
                            fetch('/admin/permissions').then(r => r.json()),
                            fetch(`/admin/roles/${roleId}/permissions`).then(r => r.json())
                        ]);

                        const allPerms = allR.permissions;
                        const roleIDs = roleR.permissions.map(p => p.id);

                        let html = '<h5>Разрешения должности:</h5>';
                        const groupedPermissions = allPerms.reduce((acc, permission) => {
                            if (!acc[permission.category]) {
                                acc[permission.category] = [];
                            }
                            acc[permission.category].push(permission);
                            return acc;
                        }, {});

                        for (const category in groupedPermissions) {
                            html += `
                                <div class="permission-category-card">
                                    <div class="permission-category-header">
                                        <span>${category}</span>
                                        <svg class="arrow-icon" viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6z"/></svg>
                                    </div>
                                    <ul class="permission-list hidden">
                            `;
                            groupedPermissions[category].forEach(p => {
                                html += `
                                    <li class="permission-item">
                                        <input class="ui-checkbox" type="checkbox" id="perm-${p.id}" data-id="${p.id}" ${roleIDs.includes(p.id) ? 'checked' : ''}>
                                        <label for="perm-${p.id}">${p.description}</label>
                                    </li>
                                `;
                            });
                            html += `</ul></div>`;
                        }

                        document.getElementById('role-permissions').innerHTML = html;

                        // Обработка кликов по заголовкам категорий
                        document.querySelectorAll('.permission-category-header').forEach(header => {
                            header.addEventListener('click', () => {
                                const list = header.nextElementSibling;
                                list.classList.toggle('hidden');
                                header.querySelector('.arrow-icon').classList.toggle('rotated');
                            });
                        });
                    } catch (error) {
                        console.error('Ошибка загрузки разрешений:', error);
                        document.getElementById('role-permissions').innerHTML =
                            '<p class="error-message">Ошибка загрузки разрешений</p>';
                    }
                });
            } catch (error) {
                console.error('Ошибка загрузки должностей:', error);
                managementContent.innerHTML =
                    '<p class="error-message">Ошибка загрузки списка должностей</p>';
            }
        });

        // --- РЕЖИМ ПОЛЬЗОВАТЕЛЕЙ ---
        document.getElementById('users-mode-button').addEventListener('click', async () => {
            document.getElementById('save-role-permissions-button').classList.add('hidden');
            document.getElementById('save-user-permissions-button').classList.remove('hidden');

            managementContent.classList.remove('hidden');
            managementContent.innerHTML = '<p>Загрузка пользователей...</p>';

            try {
                const respUsers = await fetch('/admin/users');
                if (!respUsers.ok) throw new Error('Ошибка загрузки пользователей');

                const { users } = await respUsers.json();
                const opts = users.map(u => `<option value="${u.id}">${u.full_name}</option>`).join('');
                managementContent.innerHTML = `
                    <select id="user-select" class="glass-select">
                        <option value="" disabled selected>Выберите сотрудника</option>
                        ${opts}
                    </select>
                    <div id="user-permissions" class="permissions-list"></div>
                `;

                const userSelect = document.getElementById('user-select');
                userSelect.addEventListener('change', async () => {
                    const userId = userSelect.value;
                    if (!userId) return;

                    try {
                        const [roleResponse, allU, userU] = await Promise.all([
                            fetch(`/admin/users/${userId}/role`),
                            fetch('/admin/permissions'),
                            fetch(`/admin/users/${userId}/permissions`)
                        ]);

                        if (!roleResponse.ok || !allU.ok || !userU.ok) {
                            throw new Error('Ошибка загрузки данных');
                        }

                        const role = await roleResponse.json();
                        const roleId = role.id;
                        const rolePermissions = await fetch(`/admin/roles/${roleId}/permissions`).then(r => r.json());

                        const allPerms = (await allU.json()).permissions;
                        const userIDs = (await userU.json()).userPermissions || [];
                        const rolePerms = rolePermissions.permissions.map(p => p.id);

                        let html = '<h5>Разрешения пользователя:</h5>';
                        const groupedPermissions = allPerms.reduce((acc, permission) => {
                            if (!acc[permission.category]) {
                                acc[permission.category] = [];
                            }
                            acc[permission.category].push(permission);
                            return acc;
                        }, {});

                        for (const category in groupedPermissions) {
                            html += `
                                <div class="permission-category-card">
                                    <div class="permission-category-header">
                                        <span>${category}</span>
                                        <svg class="arrow-icon" viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6z"/></svg>
                                    </div>
                                    <ul class="permission-list hidden">
                            `;
                            groupedPermissions[category].forEach(p => {
                                const isDisabled = rolePerms.includes(p.id);
                                html += `
                                    <li class="permission-item">
                                        <input class="ui-checkbox" type="checkbox" id="perm-${p.id}" data-id="${p.id}" 
                                            ${userIDs.includes(p.id) ? 'checked' : ''} 
                                            ${isDisabled ? 'disabled' : ''}>
                                        <label for="perm-${p.id}" class="${isDisabled ? 'disabled-permission' : ''}">
                                            ${p.description}
                                        </label>
                                    </li>
                                `;
                            });
                            html += `</ul></div>`;
                        }

                        document.getElementById('user-permissions').innerHTML = html;

                        // Обработка кликов по заголовкам категорий
                        document.querySelectorAll('.permission-category-header').forEach(header => {
                            header.addEventListener('click', () => {
                                const list = header.nextElementSibling;
                                list.classList.toggle('hidden');
                                header.querySelector('.arrow-icon').classList.toggle('rotated');
                            });
                        });
                    } catch (error) {
                        console.error('Ошибка загрузки разрешений:', error);
                        document.getElementById('user-permissions').innerHTML =
                            '<p class="error-message">Ошибка загрузки разрешений</p>';
                    }
                });
            } catch (error) {
                console.error('Ошибка загрузки пользователей:', error);
                managementContent.innerHTML =
                    '<p class="error-message">Ошибка загрузки списка пользователей</p>';
            }
        });
    });
}

// Блок операций с паролем (JWT)
document.getElementById('logout-button').addEventListener('click', function() {
    Swal.fire({
        title: 'Вы уверены?',
        text: 'Вы действительно хотите выйти из системы?',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Да, выйти',
        cancelButtonText: 'Отмена',
        customClass: {
            confirmButton: 'custom-button',
            cancelButton: 'custom-button'
        }
    }).then((result) => {
        if (result.isConfirmed) {
            fetch('/logout', { method: 'GET' })
                .then(response => {
                    if (response.ok) {
                        window.location.href = '/';
                    } else {
                        Swal.fire('Ошибка!', 'Не удалось выйти из системы.', 'error');
                    }
                })
                .catch(() => Swal.fire('Ошибка!', 'Произошла ошибка при выходе.', 'error'));
        }
    });
});
document.getElementById('change-password-button').addEventListener('click', function() {
    const username = main.dataset.username;

    showModal('Сменить пароль', `
        <form id="changePasswordForm">
            <div class="form-group">
                <label for="currentPassword">Текущий пароль:</label>
                <input type="text" id="currentPassword" class="input-field" placeholder="Введите текущий пароль">
            </div>
            <div class="form-group">
                <label for="newPassword">Новый пароль:</label>
                <input type="text" id="newPassword" class="input-field" placeholder="Введите новый пароль">
            </div>
            <div class="form-group">
                <label for="confirmPassword">Подтвердите новый пароль:</label>
                <input type="text" id="confirmPassword" class="input-field" placeholder="Подтвердите новый пароль">
            </div>
        </form>
    `, () => {
        const currentPassword = document.getElementById('currentPassword').value.trim();
        const newPassword = document.getElementById('newPassword').value.trim();
        const confirmPassword = document.getElementById('confirmPassword').value.trim();

        if (!currentPassword || !newPassword || !confirmPassword) {
            Swal.showValidationMessage('Пожалуйста, заполните все поля.');
            return false;
        }

        if (newPassword !== confirmPassword) {
            Swal.showValidationMessage('Новый пароль и его подтверждение не совпадают.');
            return false;
        }

        // Проверка на длину и сложность пароля
        const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$/;
        if (!passwordRegex.test(newPassword)) {
            Swal.showValidationMessage('Пароль должен содержать минимум 8 символов, включая буквы, цифры и специальные знаки.');
            return false;
        }

        return fetch('/change-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, currentPassword, newPassword })
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    return data;
                } else {
                    Swal.showValidationMessage(data.error || 'Ошибка при смене пароля.');
                }
            })
            .catch(() => Swal.showValidationMessage('Ошибка запроса к серверу.'));
    }, '600px').then(result => {
        if (result.isConfirmed) {
            Swal.fire({
                title: "Изменено!",
                text: "Пароль успешно изменен.",
                icon: "success",
                timer: 1000,
                timerProgressBar: true,
                toast: true,
                position: "top-end",
                showConfirmButton: false
            }).then(() => {
                window.location.reload();
            });

        }
    });
});
