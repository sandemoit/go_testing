package routes

import "github.com/gofiber/fiber/v2"

func PublicRoute(route *fiber.App) {
	RoleRoute(route)
	UserRoute(route)
	ProdRoute(route)
	EfisiensiRoute(route)
}
