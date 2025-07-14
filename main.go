package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	approuter "IoTT/app/router"
	"IoTT/internal/database"
	"IoTT/internal/models"
	internalrouter "IoTT/internal/router"
	"IoTT/internal/telegram"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Gagal memuat file .env. Menggunakan environment variable sistem.")
	}

	database.InitDB()
	if database.DB != nil {
		defer database.CloseDB()
	}

	telegram.LoadConfig()
	if err := telegram.InitBot(); err != nil {
		log.Printf("Peringatan: Gagal menginisialisasi bot Telegram: %v. Notifikasi mungkin tidak berfungsi.", err)
	}

	models.StartOfflineDetectionWorker()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", // Untuk produksi, ganti dengan domain frontend Anda, misal: "http://localhost:3000,https://app.domain.com"
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	internalrouter.SetupInternalRouter(app)

	appAPIGroup := app.Group("/api/v1")
	approuter.SetupAppRouter(appAPIGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("🚀 Memulai server Fiber di port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server Fiber: %v", err)
	}
}
