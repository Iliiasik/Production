package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"production/database"
	"production/models"
	"time"
)

func ShowLoginPage(c *gin.Context) {
	// Пытаемся получить токен из куки
	tokenStr, err := c.Cookie("token")
	if err == nil {
		// Проверяем валидность токена
		claims, err := ValidateJWT(tokenStr)
		if err == nil && claims != nil {
			// Токен валиден, перенаправляем на /home
			c.Redirect(302, "/home")
			return
		}
	}

	// Если нет токена или он невалиден, показываем страницу логина
	c.HTML(200, "login.html", gin.H{})
}
func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("неверные claims")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 мощный хеш, но долгий
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Секретный ключ, используемый для подписания и валидации JWT (Для прода необходимо засунуть его в env)
var JwtKey = []byte("secret_key")

type Claims struct {
	EmployeeID uint   `json:"employee_id"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(employeeID uint, username string) (string, error) {
	// Создает JWT-токен, содержащий employeeID и username.
	//
	//ExpiresAt — токен истекает через 24 часа.
	//
	//IssuedAt — время выпуска токена.
	//
	//Используется алгоритм HS256 и секретный ключ JwtKey.
	//
	//Возвращает подписанный токен в виде строки.
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		EmployeeID: employeeID,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var employee models.Employee
	if err := database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		c.HTML(401, "login.html", gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if !CheckPasswordHash(password, employee.PasswordHash) {
		c.HTML(401, "login.html", gin.H{"error": "Неверный логин или пароль"})
		return
	}

	token, err := GenerateJWT(employee.ID, employee.Username)
	if err != nil {
		c.HTML(500, "login.html", gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Сохраняем токен в куку
	c.SetCookie("token", token, 3600*24, "/", "localhost", false, true)

	// Перенаправляем на личный кабинет
	c.Redirect(302, "/home")
}

func ChangePassword(c *gin.Context) {
	var request struct {
		Username        string `json:"username"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	// Прочитать тело запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Находим сотрудника по логину
	var employee models.Employee
	if err := database.DB.Where("username = ?", request.Username).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин"})
		return
	}

	// Проверяем текущий пароль
	if !CheckPasswordHash(request.CurrentPassword, employee.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный текущий пароль"})
		return
	}

	// Хешируем новый пароль
	hashedPassword, err := HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	// Обновляем запись сотрудника в базе данных
	employee.PasswordHash = hashedPassword
	employee.IsPasswordChanged = true // Устанавливаем флаг, что пароль был изменен
	if err := database.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пароля"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Пароль успешно изменен"})
}

func Logout(c *gin.Context) {
	// Удаляем токен из cookies
	c.SetCookie("token", "", -1, "/", "localhost", false, true) // Устанавливаем токен с истекшим временем

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Вы успешно вышли из системы"})
}
