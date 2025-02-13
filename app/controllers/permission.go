package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"

	"github.com/gofiber/fiber/v2"
)

func CreatePermission(c *fiber.Ctx) error {
	var body struct {
		Permission string `json:"permission"`
		RoleID     uint   `json:"role_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid request payload",
		})
	}

	// Cari roles berdasarkan RoleIDs
	var roles []models.Role
	database.DB.Find(&roles, body.RoleID)

	// Buat permission tanpa mengisi Id
	permission := models.Permission{
		Permission: body.Permission,
		Roles:      roles,
	}

	// Simpan ke database
	result := database.DB.Create(&permission)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to create: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Successfully created " + permission.Permission,
	})
}

func DeletePermission(c *fiber.Ctx) error {
	permission := c.Params("permission")

	database.DB.Delete(&models.Permission{}, "permission = ?", permission)

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Successfully deleted " + permission,
	})
}

func FilterPermission(c *fiber.Ctx) error {
	permissionQuery := c.Query("permission")
	if permissionQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Permission parameter is required",
		})
	}

	var permissions []models.Permission
	result := database.DB.
		Preload("Roles").       // Memuat permission dalam roles
		Preload("Roles.Users"). // Memuat users dalam roles
		Where("permission LIKE ?", "%"+permissionQuery+"%").
		Find(&permissions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to filter permissions: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}

func GetPermission(c *fiber.Ctx) error {
	permissionQuery := c.Query("permission")

	var permissions []models.Permission
	result := database.DB.
		Preload("Roles").       // Memuat permission dalam roles
		Preload("Roles.Users"). // Memuat users dalam roles
		Where("permission LIKE ?", "%"+permissionQuery+"%").
		Find(&permissions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to filter permissions: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}
