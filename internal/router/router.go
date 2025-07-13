package router

import (
	"github.com/gofiber/fiber/v2" // GANTI: Menggunakan Fiber

	"IoTT/internal/handler"
	// "IoTT/internal/telegram" // Kita akan gunakan handler dari paket 'handler'
)

// SetupInternalRouter mengkonfigurasi rute untuk layanan internal menggunakan Fiber.
// GANTI: Parameter diubah dari *gin.Engine menjadi *fiber.App
func SetupInternalRouter(app *fiber.App) {

	// Grup rute untuk data sensor
	// Konsepnya sama, hanya dipanggil dari `app` (Fiber) bukan `router` (Gin)
	dataRoutes := app.Group("/data")
	{
		// GANTI: Method .POST() menjadi .Post() (huruf kecil)
		dataRoutes.Post("", handler.HandleSensorData)
	}

	telegramRoutes := app.Group("/telegram")
	{
		// GANTI: Menggunakan handler yang sudah kita migrasikan dari paket 'handler'
		// dan method .POST() menjadi .Post()
		telegramRoutes.Post("/webhook", handler.HandleTelegramWebhook)
	}
}
