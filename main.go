package main

import (
	"log"
	"os"

	// Asumsikan path ini masih sama, sesuaikan jika perlu

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
	// 1. Muat konfigurasi dari file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Gagal memuat file .env. Menggunakan environment variable sistem.")
	}

	// 2. Inisialisasi dependensi utama
	database.InitDB()
	if database.DB != nil {
		defer database.CloseDB() // Gunakan defer untuk menutup koneksi saat main() selesai
	}

	// Logika Telegram tetap sama
	telegram.LoadConfig()
	if err := telegram.InitBot(); err != nil {
		log.Printf("Peringatan: Gagal menginisialisasi bot Telegram: %v. Notifikasi mungkin tidak berfungsi.", err)
	}

	// 3. Memulai proses background (worker) - tidak ada perubahan
	models.StartOfflineDetectionWorker()

	// 4. Setup server Fiber
	app := fiber.New()

	// 4a. Tambahkan Middleware CORS (versi Fiber)
	// Konfigurasi ini sangat mirip dengan versi Gin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", // Untuk produksi, ganti dengan domain frontend Anda, misal: "http://localhost:3000,https://app.domain.com"
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	// 5. Daftarkan rute
	// PERHATIAN: Anda perlu mengubah isi dari fungsi SetupInternalRouter dan SetupAppRouter
	// agar menerima parameter dari Fiber, bukan Gin.
	internalrouter.SetupInternalRouter(app) // Sekarang menerima `*fiber.App`

	appAPIGroup := app.Group("/api/v1")   // Membuat grup rute di Fiber
	approuter.SetupAppRouter(appAPIGroup) // Sekarang menerima `fiber.Router`

	// 6. Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Fallback port
	}
	log.Printf("🚀 Memulai server Fiber di port %s", port)

	// Gunakan app.Listen() untuk Fiber
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server Fiber: %v", err)
	}
}
