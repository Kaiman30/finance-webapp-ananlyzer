package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

// Инициализация БД
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./finance.db")
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	// Автоматическая миграция для создания таблиц
	if err := db.AutoMigrate(&Transactions{}, &User{}).Error; err != nil {
		return nil, fmt.Errorf("не удалось создать таблицы: %v", err)
	}

	DB = db
	return db, nil
}

// Модель для транзакции
type Transactions struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
}

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
