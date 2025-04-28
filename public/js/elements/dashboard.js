const main = document.getElementById('user-dashboard');
const fullName = main.dataset.fullname;
const position = main.dataset.position;
const salary = main.dataset.salary;

const html = `
<section class="cards-wrapper">
  <div class="card card--large">
    <div class="profile-info">
      <h2>Добро пожаловать, ${fullName}!</h2>

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
    </div>
    </div>
  </div>
</section>
`;

main.innerHTML = html;


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
                <button id="roles-mode-button" class="profile-button">Роли</button>
                <button id="users-mode-button" class="profile-button">Пользователи</button>
            </div>
            <div id="management-content" class="management-content hidden"></div>
            <div class="actions-center">
                <button id="close-permissions" class="profile-button">Закрыть</button>
                <button id="save-role-permissions-button" class="profile-button hidden">Сохранить роли</button>
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
                    title: 'Выберите роль',
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
                    title: 'Роль обновлена',
                    showConfirmButton: false,
                    timer: 3000
                });
            } catch (error) {
                console.error('Ошибка сохранения роли:', error);
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
                if (!respRoles.ok) throw new Error('Ошибка загрузки ролей');

                const { positions } = await respRoles.json();
                const opts = positions.map(r => `<option value="${r.id}">${r.name}</option>`).join('');
                managementContent.innerHTML = `
                    <select id="role-select" class="glass-select">
                        <option value="" disabled selected>Выберите роль</option>
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

                        let html = '<h5>Разрешения роли:</h5>';
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
                console.error('Ошибка загрузки ролей:', error);
                managementContent.innerHTML =
                    '<p class="error-message">Ошибка загрузки списка ролей</p>';
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
