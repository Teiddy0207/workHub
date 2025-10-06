package main

import (
	"fmt"
	"log"
	"os"
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

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	fmt.Printf("🚀 Server đang chạy tại: http://localhost:%s\n", port)
	log.Println("🚀 Server đang chạy ở cổng:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
