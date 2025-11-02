package main

import (
	"fmt"
	"log"
	"os"
	"workHub/config"
	"workHub/logger"
	"workHub/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ Không tìm thấy .env, đang dùng mặc định")
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf("❌ Không thể tạo thư mục logs: %v", err)
		}
	}

	logger.InitLogger()
	defer logger.Log.Sync()

	// Kết nối database và migrate
	logger.Info("main", "main", "Connecting to database...")
	db, err := config.ConnectDatabase()
	if err != nil {
		logger.Error("main", "main", fmt.Sprintf("Failed to connect database: %v", err))
		log.Fatal(err)
	}
	logger.Info("main", "main", "Database connected successfully")

	// Auto migrate tables
	if err := config.AutoMigrate(db); err != nil {
		logger.Error("main", "main", fmt.Sprintf("Failed to migrate database: %v", err))
		log.Fatal(err)
	}

	r := router.InitRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"  
	}

	logger.Info("main", "main", "Server đang chạy tại: http://localhost:"+port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
