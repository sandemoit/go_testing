package routes

import (
	"go_teknologi/utils"

	_ "go_teknologi/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(route *fiber.App) {
	utils.LogAPI(route)

	// Middleware CORS
	route.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Ganti dengan domain frontend kamu
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false, // Jika pakai cookies/token di header
	}))

	// Contoh route untuk testing
	route.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Frontend!",
		})
	})

	// Route Swagger UI
	route.Get("/swagger/*", fiberSwagger.WrapHandler)

	PublicRoute(route)
}
