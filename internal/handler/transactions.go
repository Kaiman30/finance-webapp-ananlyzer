package handler

import (
	"finance-app/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

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
		if err := c.ShouldBindBodyWithJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат данных"})
			return
		}

		// Валидация данных
		if errors := validateTransaction(&transaction); len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		if err := db.Create(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить транзакцию"})
			return
		}

		c.JSON(http.StatusCreated, transaction) // Возвращаем добавленную транзакцию
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
