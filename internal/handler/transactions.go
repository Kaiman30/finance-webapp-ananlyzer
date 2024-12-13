package handler

import (
	"finance-app/internal/storage"
	"net/http"

	//"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//"golang.org/x/crypto/bcrypt"
)

//var jwtSecret = []byte("Test12345")

// Защита маршрутов с помощью JWT
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
// 			c.Abort()
// 			return
// 		}

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return jwtSecret, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверый токен"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()

// 	}
// }

// Регистрация пользователя
// func Register(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user storage.User
// 		if err := c.ShouldBindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
// 			return
// 		}

// 		// Хешируем пароль
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибки при хешировании пароля"})
// 			return
// 		}
// 		user.Password = string(hashedPassword)

// 		// Сохраняем пользователя
// 		if err := db.Create(&user).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить пользователя"})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, gin.H{"message": "Пользователь зарегистрирован"})
// 	}
// }

// // Вход пользователя
// func Login(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user storage.User
// 		var loginUser storage.User
// 		if err := c.ShouldBindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
// 			return
// 		}

// 		// Проверяем, существует ли пользователь
// 		if err := db.Where("username = ?", user.Username).First(&loginUser).Error; err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
// 			return
// 		}

// 		// Проверяем пароль
// 		if err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password)); err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
// 			return
// 		}

// 		// Создаем JWT токен
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"username": user.Username,
// 			"exp":      time.Now().Add(time.Hour * 24).Unix(),
// 		})

// 		tokenString, err := token.SignedString(jwtSecret)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании JWT-токена"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"token": tokenString})
// 	}
// }

// Валидация транзакции
func validateTransaction(transaction *storage.Transactions) []string {
	var errors []string

	// Проверка на обязательность суммы и корректность
	if transaction.Amount <= 0 {
		errors = append(errors, "Сумма должна быть больше нуля")
	}

	// Проверка на корректность описания
	if len(transaction.Description) < 5 {
		errors = append(errors, "Описание должно иметь минимум 5 символов")
	}

	// Проверка на корректность категории
	validCategories := []string{"Еда", "Проезд", "Развлечения", "Долги", "Жилье", "Подписки"}
	categoryValid := false
	for _, category := range validCategories {
		if category == transaction.Category {
			categoryValid = true
			break
		}
	}
	if !categoryValid {
		errors = append(errors, "Некорректная категория")
	}

	return errors
}

// Получение всех транзакций
func GetTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactions []storage.Transactions
		if err := db.Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить транзакции"})
			return
		}
		c.JSON(http.StatusOK, transactions)
	}
}

// Добавление новой транзакции
func AddTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction storage.Transactions
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		// Извлекаем ID пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			return
		}

		// Добавляем ID пользователя в транзакцию
		transaction.ID = userID.(uint)

		// Сохраняем транзакцию в базе данных
		if err := db.Create(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении транзакции"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Транзакция добавлена"})
	}
}

func UpdateTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction storage.Transactions
		id := c.Param("id")

		// Ищем транзакцию по ID
		if err := db.First(&transaction, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Транзакция не найдена"})
			return
		}

		// Обновляем данные транзакции
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат данных"})
			return
		}

		// Валидация данных
		if errors := validateTransaction(&transaction); len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		// Сохраняем обновленную транзакцию
		if err := db.Save(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить транзакцию"})
			return
		}

		c.JSON(http.StatusOK, transaction)
	}
}

// Удаление транзакции
func DeleteTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction storage.Transactions
		id := c.Param("id")

		// Ищем транзакцию по ID
		if err := db.First(&transaction, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Транзакция не найдена"})
			return
		}

		// Удаляем транзакцию
		if err := db.Create(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить транзакцию"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Транзакция удалена"})
	}
}
