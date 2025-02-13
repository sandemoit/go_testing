package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	var body struct {
		Role         string `json:"role"`
		PermissionID uint   `json:"permission_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid request payload",
		})
	}

	// Cari roles berdasarkan RoleIDs
	var roles []models.Role
	database.DB.Find(&roles, body.PermissionID)

	// Buat permission tanpa mengisi Id
	role := models.Role{
		Role: body.Role,
	}

	// Simpan ke database
	result := database.DB.Create(&role)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to create role",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Successfully created role",
	})
}

func DeleteRole(c *fiber.Ctx) error {
	role := c.Params("role")

	database.DB.Delete(&models.Role{}, "role = ?", role)

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Successfully deleted role",
	})
}

func FilterRole(c *fiber.Ctx) error {
	roleQuery := c.Query("role")
	if roleQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Role parameter is required",
		})
	}

	var roles []models.Role
	result := database.DB.
		Preload("Users").             // Memuat users dalam roles
		Preload("Users.Roles").       // Memuat roles dalam users
		Preload("Permissions").       // Memuat permissions dalam roles
		Preload("Permissions.Roles"). // Memuat roles dalam permissions
		Where("role LIKE ?", "%"+roleQuery+"%").
		Find(&roles)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to filter roles: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(roles)
}

func GetRole(c *fiber.Ctx) error {
	roleQuery := c.Query("role")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	offset := (page - 1) * limit

	var roles []models.Role
	var total int64

	// Hitung total data
	database.DB.Model(&models.Role{}).Where("role LIKE ?", "%"+roleQuery+"%").Count(&total)

	// Query dengan pagination
	result := database.DB.
		Preload("Users").             // Memuat users dalam roles
		Preload("Users.Roles").       // Memuat roles dalam users
		Preload("Permissions").       // Memuat permissions dalam roles
		Preload("Permissions.Roles"). // Memuat roles dalam permissions
		Where("role LIKE ?", "%"+roleQuery+"%").
		Offset(offset).
		Limit(limit).
		Find(&roles)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to filter roles: " + result.Error.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status: true,
		Data:   roles,
		Message: fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
