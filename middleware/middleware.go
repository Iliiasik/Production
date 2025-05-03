package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"production/controllers"
	"production/database"
	"production/models"
)

// Authorize проверяет, есть ли у текущего пользователя разрешение permissionName.
// Токен берётся из cookie "token", в которой сохраняется JWT после входа.
func Authorize(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			showErrorPage(c, "Пожалуйста, авторизуйтесь")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &controllers.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})
		if err != nil || !token.Valid {
			showErrorPage(c, "Неверный или просроченный токен")
			return
		}

		claims, ok := token.Claims.(*controllers.Claims)
		if !ok {
			showErrorPage(c, "Неверные данные токена")
			return
		}
		employeeID := claims.EmployeeID

		var perm models.Permission
		if err := database.DB.Where("name = ?", permissionName).First(&perm).Error; err != nil {
			showErrorPage(c, "Разрешение не найдено")
			return
		}

		var userPermissions []models.UserPermission
		if err := database.DB.
			Where("employee_id = ?", employeeID).
			Find(&userPermissions).Error; err != nil {
			showErrorPage(c, "Ошибка загрузки прав пользователя")
			return
		}

		var userHasPermission bool
		for _, up := range userPermissions {
			if up.PermissionID == perm.ID {
				userHasPermission = true
				break
			}
		}
		if userHasPermission {
			c.Next()
			return
		}

		var emp models.Employee
		if err := database.DB.First(&emp, employeeID).Error; err != nil {
			showErrorPage(c, "Пользователь не найден")
			return
		}

		var positionPermissions []models.PositionPermission
		if err := database.DB.
			Where("position_id = ? AND permission_id = ?", emp.PositionID, perm.ID).
			Find(&positionPermissions).Error; err == nil && len(positionPermissions) > 0 {
			c.Next()
			return
		}

		showErrorPage(c, "У вас нет прав для выполнения действия")
	}
}

func showErrorPage(c *gin.Context, message string) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusForbidden, `
<html>
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="/assets/css/elements/sweetalert-custom.css">
    <link rel="stylesheet" href="/assets/css/general/layout.css">
    <script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
</head>
<body>
    <script>
        Swal.fire({
            icon: 'error',
            title: 'Ошибка доступа',
            text: '%s',
            confirmButtonText: 'Вернуться назад',
            customClass: {
                popup: 'popup-class',
                confirmButton: 'custom-button',
                cancelButton: 'custom-button'
            }
        }).then(() => {
            window.history.back();
        });
    </script>
</body>
</html>
	`, message)
	c.Abort()
}
