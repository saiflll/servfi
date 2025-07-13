package router

import (
	"github.com/gofiber/fiber/v2" // GANTI: Menggunakan Fiber

	"IoTT/app/handler"
	"IoTT/app/middleware"
)

// SetupAppRouter mengkonfigurasi rute aplikasi utama menggunakan Fiber.
// GANTI: Parameter diubah dari *gin.RouterGroup menjadi fiber.Router
func SetupAppRouter(router fiber.Router) {
	// Rute Otentikasi
	authGroup := router.Group("/auth")
	{
		// GANTI: .POST menjadi .Post (huruf kecil)
		authGroup.Post("/login", handler.Login)
		// Middleware dipanggil dengan cara yang sama
		authGroup.Post("/logout", middleware.AuthMiddleware(), handler.Logout)
	}

	// Grup rute yang dilindungi oleh middleware
	protected := router.Group("/")
	// GANTI: Cara menggunakan middleware sama, tapi middleware-nya sendiri sudah versi Fiber
	protected.Use(middleware.AuthMiddleware())
	{
		// GANTI: Semua method diubah ke huruf kecil (e.g., GET -> Get)
		// Sintaks parameter (:area_id) tetap sama.
		protected.Get("/area/:area_id/:no/summary", handler.GetAreaSummaryBySensorHandler)
		protected.Get("/area/:area_id/sensors/export", handler.ExportSensorDataToExcel)
		protected.Get("/datas", handler.GetSensorStatuses)
		protected.Get("/alerts/detailed", handler.GetDetailedAlerts)
		protected.Get("/chart-data", handler.GetChartDataHandler)
	}
}
