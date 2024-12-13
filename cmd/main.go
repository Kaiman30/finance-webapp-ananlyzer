package main

import (
	"finance-app/internal/handler"
	"finance-app/internal/middleware"
	"finance-app/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}

	// Инициализация Gin
	r := gin.Default()

	// Маршруты для авторизации
	r.POST("/api/register", handler.Register(db))
	r.POST("/api/login", handler.Login(db))

	// Маршруты для работы с транзакциями
	r.GET("/api/transactions", middleware.AuthMiddleware(), handler.GetTransactions(db))
	r.POST("/api/transactions", middleware.AuthMiddleware(), handler.AddTransactions(db))
	r.PUT("/api/transactions", middleware.AuthMiddleware(), handler.UpdateTransactions(db))
	r.DELETE("/api/transactions", middleware.AuthMiddleware(), handler.DeleteTransactions(db))

	// Отдача статических файлов (правильный путь для /web)
	r.Static("/web", "./web") // Все файлы из папки ./web будут доступны по пути /web

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
