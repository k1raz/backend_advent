package auth

import (
	"backend/internal/db"
	"backend/types"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func generateToken(userId int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Register(c *gin.Context) {
	var json types.RegisterPayload

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
		log.Printf("DB Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "invalid",
			"details": "user_creation_failed",
		})
		return
	}

	userId, err := result.LastInsertId()
	if err != nil {
		log.Printf("LastInsertId Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "invalid",
			"details": "id_retrieval_failed",
		})
		return
	}

	token, err := generateToken(int(userId), json.Username)
	if err != nil {
		log.Printf("Token Generation Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "invalid",
			"details": "token_generation_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
}

func Login(c *gin.Context) {
	var json types.LoginPayload

	fmt.Println(json)

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение пользователя из базы данных
	var userId int
	var passwordHash string
	err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?",
		json.Username).Scan(&userId, &passwordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(json.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	token, err := generateToken(userId, json.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешный вход", "token": token})
}

func ResetPassword(c *gin.Context) {
	var json types.RestorePassword

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userId int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = ?", json.Username).Scan(&userId)

	fmt.Printf("Username: %s, UserId: %d, Error: %v\n", json.Username, userId, err)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске пользователя"})
		}
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(json.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	result, err := db.DB.Exec("UPDATE users SET password_hash = ? WHERE username = ?",
		passwordHash, json.Username)

	if err != nil {
		fmt.Printf("Update Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пароля"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пароль"})
		return
	}

	token, err := generateToken(userId, json.Username)
	if err != nil {
		log.Printf("Token Generation Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "invalid",
			"details": "token_generation_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
}
