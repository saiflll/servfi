package router

import (
	"github.com/gofiber/fiber/v2" 

	"IoTT/app/handler"
	"IoTT/app/middleware"
)



func SetupAppRouter(router fiber.Router) {
	
	authGroup := router.Group("/auth")
	{
		
		authGroup.Post("/login", handler.Login)
		
		authGroup.Post("/logout", middleware.AuthMiddleware(), handler.Logout)
	}

	
	protected := router.Group("/")
	
	protected.Use(middleware.AuthMiddleware())
	{
		
		
		protected.Get("/area/:area_id/:no/summary", handler.GetAreaSummaryBySensorHandler)
		protected.Get("/area/:area_id/sensors/export", handler.ExportSensorDataToExcel)
		protected.Get("/export/excel", handler.ExportSingleSensorDataToExcel)
		protected.Get("/datas", handler.GetSensorStatuses)
		protected.Get("/alerts/detailed", handler.GetDetailedAlerts)
		protected.Get("/chart-data", handler.GetChartDataHandler)
	}
}
