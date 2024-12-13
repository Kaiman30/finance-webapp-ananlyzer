package handler

import (
	"finance-app/internal/config"
	"finance-app/internal/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Структура для Claims (полезная нагрузка)
type Claims struct {
	Sub uint `json:"sub"` // Идентификатор пользователя
	jwt.StandardClaims
}

// Генерация JWT токена
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		Sub: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Токен истекает через 24 часа
		},
	}

	// Здесь мы указываем алгоритм подписи через метод HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием вашего секрета
	return token.SignedString([]byte(config.JWTSecret)) // Используем секретный ключ из config
}

// Функция для парсинга и проверки токена
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Парсим токен с использованием функции, которая указывает на метод подписи
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что токен использует алгоритм HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
		}
		// Возвращаем секретный ключ для верификации
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Если токен валиден, извлекаем и возвращаем Claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("неверный токен")
}

// Регистрация пользователя
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user storage.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
			return
		}
		user.Password = string(hashedPassword)

		// Сохраняем пользователя в БД
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Пользователь зарегистрирован"})
	}
}

// Авторизация пользователя
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials storage.User
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		var user storage.User
		if err := db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Сравниваем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Создаем JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString(config.JWTSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
