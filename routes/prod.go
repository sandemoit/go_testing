package routes

import (
	"go_teknologi/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProdRoute(route *fiber.App) {
	prodGroup := route.Group("/api")

	prodGroup.Post("/import-prod_dailies", controllers.ImportExcel)
	prodGroup.Get("/getdraft-prod_dailies", controllers.GetDataDraft)
	prodGroup.Get("/getrilis-prod_dailies", controllers.GetDataRilis)
	prodGroup.Put("/update-prod_dailies", controllers.UpdateProd)
	prodGroup.Put("/revoke-prod_dailies", controllers.RevokeProd)
}
