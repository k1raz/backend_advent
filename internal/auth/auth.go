package auth

import (
	"backend/internal/db"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

func generateToken(userId int, username string) (string, error) {
	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(), // Токен действителен 7 дней
	})

	// Подписываем токен
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Register(c *gin.Context) {
    var json struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
        return
    }

    result, err := db.DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", 
	json.Username, passwordHash)
	
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid"})
        return
    }

    userId, err := result.LastInsertId()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid"})
        return
    }

    // Генерация JWT токена
    token, err := generateToken(int(userId), json.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
}

func Login(c *gin.Context) {
    var json struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Получение пользователя из базы данных
    var userId int
    var passwordHash string
    err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", json.Username).Scan(&userId, &passwordHash)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    // Проверка пароля
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(json.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    // Генерация JWT токена
    token, err := generateToken(userId, json.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Успешный вход", "token": token})
} 