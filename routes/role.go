package routes

import (
	"go_teknologi/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func RoleRoute(route *fiber.App) {
	roleGroup := route.Group("/api")

	roleGroup.Post("/create-permission", controllers.CreatePermission)
	roleGroup.Delete("/delete-permission/:permission", controllers.DeletePermission)
	roleGroup.Get("/filter-permission", controllers.FilterPermission)
	roleGroup.Get("/getdata-permission", controllers.GetPermission)

	roleGroup.Post("/create-role", controllers.CreateRole)
	roleGroup.Delete("/delete-role/:role", controllers.DeleteRole)
	roleGroup.Get("/filter-role", controllers.FilterRole)
	roleGroup.Get("/getdata-role", controllers.GetRole)

	roleGroup.Post("/create-role_permission", controllers.CreateRolePermission)
	roleGroup.Delete("/delete-role_permission/:id", controllers.DeleteRolePermission)
	roleGroup.Get("/get-role_permission", controllers.GetRolePermission)
	roleGroup.Get("/get-role_permission/:role_id", controllers.GetRolePermissionById)

}
