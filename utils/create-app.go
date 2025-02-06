package utils

import "github.com/gofiber/fiber/v2"

func CreateApp() *fiber.App {

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100 mb limit uploud files
	})

	return app
}
