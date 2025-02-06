package routes

import (
	"go_teknologi/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route *fiber.App) {
	userGroup := route.Group("/api")

	userGroup.Post("/create-user", controllers.CreateUser)
	userGroup.Delete("/delete-user/:username", controllers.DeleteUser)
	userGroup.Get("/filter-user", controllers.FilterUser)
	userGroup.Get("/getdata-user", controllers.GetDataUser)

	userGroup.Post("/create-user_role", controllers.CreateUserRole)
	userGroup.Delete("/delete-user_role/:id", controllers.DeleteUserRole)
	userGroup.Get("/get-user_role", controllers.GetUserRole)
	userGroup.Get("/get-user_role/:user_id", controllers.GetUserRoleById)
}
