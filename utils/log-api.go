package utils

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LogAPI(route *fiber.App) error {

	// Middleware to log API access
	route.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		// Process the request
		err := c.Next()
		duration := time.Since(start)

		// Log method, path, and status code
		log.Printf("%s %s %d - %s",
			// time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	})

	return nil
}
