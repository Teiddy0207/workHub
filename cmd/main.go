package main

import (
	"fmt"
	"log"
	"os"
	"workHub/config"
	internalconfig "workHub/internal/config"
	"workHub/logger"
	"workHub/router"

	"workHub/pkg/utils"

	"github.com/joho/godotenv"
)

func InitApp(cfg *internalconfig.Config) {
	if cfg.Redis != nil {
		utils.InitRedis(
			cfg.Redis.Host,
			cfg.Redis.Port,
			cfg.Redis.Password,
			cfg.Redis.DB,
			cfg.Redis.PoolSize,
			cfg.Redis.ReadTimeout,
			cfg.Redis.WriteTimeout,
			cfg.Redis.DialTimeout,
			cfg.Redis.Timeout,
		)
	}
}

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

	cfg, err := internalconfig.LoadConfig()
	if err != nil {
		logger.Error("main", "main", fmt.Sprintf("Failed to load config: %v", err))
		log.Fatal(err)
	}

	// Khởi tạo Redis
	InitApp(cfg)

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
