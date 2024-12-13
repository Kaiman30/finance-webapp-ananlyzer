package main

import (
	"finance-app/internal/handler"
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

	// Маршруты для API
	r.GET("/api/transactions", handler.GetTransactions(db))
	r.POST("/api/transactions", handler.AddTransactions(db))
	r.PUT("/api/transactions", handler.UpdateTransactions(db))
	r.DELETE("/api/transactions", handler.DeleteTransactions(db))

	// Отдача статических файлов (правильный путь для /web)
	r.Static("/web", "./web") // Все файлы из папки ./web будут доступны по пути /web

	// Для главной страницы, отдаём index.html при запросе по /web
	// r.GET("/web", func(c *gin.Context) {
	// 	c.File("./web/index.html")
	// })

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
