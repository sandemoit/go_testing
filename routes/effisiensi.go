package routes

import (
	"go_teknologi/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func EfisiensiRoute(route *fiber.App) {
	efisiensiGroup := route.Group("/api")

	efisiensiGroup.Get("/getdata-effisiensi_amoniak", controllers.GetDataEfisiensiAmoniak)
	efisiensiGroup.Get("/getdata-effisiensi_urea", controllers.GetDataEfisiensiUrea)
	efisiensiGroup.Get("/getdata-effisiensi_utilitas", controllers.GetDataEfisiensiUtilitas)
}
