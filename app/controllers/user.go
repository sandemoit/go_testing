package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(ctx *fiber.Ctx) error {
	var request models.EmployeeRequest

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid request body",
		})
	}

	// Buat object user baru
	user := models.User{
		Name:     request.Name,
		Username: request.Badge,
	}

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal membuat user baru",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "User berhasil dibuat",
		// Data:    user,
	})
}

func DeleteUser(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	// Hapus user_roles terlebih dahulu
	if err := database.DB.Where("user_id IN (?)",
		database.DB.Model(&models.User{}).Select("id").Where("username = ?", username),
	).Delete(&models.UserRole{}).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal menghapus user roles",
		})
	}

	// Kemudian hapus user
	if err := database.DB.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal menghapus user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "User berhasil dihapus",
	})
}

func FilterUser(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	username := ctx.Query("username")

	var users []models.User
	query := database.DB.Preload("Roles.Permissions").Model(&models.User{})

	// Filter berdasarkan nama
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Filter berdasarkan username
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// Eksekusi query
	if err := query.Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal mendapatkan data user",
		})
	}

	// Transform data agar sesuai format JSON yang diinginkan
	var transformedUsers []map[string]interface{}

	for _, user := range users {
		var roles []map[string]interface{}
		for _, role := range user.Roles {
			var permissions []map[string]interface{}
			for _, perm := range role.Permissions {
				permissions = append(permissions, map[string]interface{}{
					"id":         perm.Id,
					"permission": perm.Permission,
					"roles":      []string{role.Role}, // Menyimpan role terkait
				})
			}
			roles = append(roles, map[string]interface{}{
				"id":          role.Id,
				"role":        role.Role,
				"permissions": permissions,
				"users":       []string{user.Name}, // Menyimpan nama user dalam role
			})
		}

		transformedUsers = append(transformedUsers, map[string]interface{}{
			"id":             user.Id,
			"name":           user.Name,
			"username":       user.Username,
			"remember_token": user.RememberToken,
			"roles":          roles,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(transformedUsers)
}

func GetDataUser(ctx *fiber.Ctx) error {
	// Ambil parameter dari query
	name := ctx.Query("name")
	username := ctx.Query("username")
	limit, _ := strconv.Atoi(ctx.Query("limit", "10")) // Default 10
	page, _ := strconv.Atoi(ctx.Query("page", "1"))    // Default 1

	var users []models.User
	query := database.DB.Preload("Roles.Permissions").Model(&models.User{})

	// Filter berdasarkan nama jika tidak kosong
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Filter berdasarkan username jika tidak kosong
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// Hitung total data sebelum pagination
	var totalRecords int64
	query.Count(&totalRecords)

	// Hitung offset untuk pagination
	offset := (page - 1) * limit

	// Eksekusi query dengan limit & offset
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal mendapatkan data user",
		})
	}

	// Transformasi data agar sesuai format JSON yang diinginkan
	var transformedUsers []map[string]interface{}

	for _, user := range users {
		var roles []map[string]interface{}
		for _, role := range user.Roles {
			var permissions []map[string]interface{}
			for _, perm := range role.Permissions {
				permissions = append(permissions, map[string]interface{}{
					"id":         perm.Id,
					"permission": perm.Permission,
					"roles":      []string{role.Role},
				})
			}
			roles = append(roles, map[string]interface{}{
				"id":          role.Id,
				"role":        role.Role,
				"permissions": permissions,
				"users":       []string{user.Name},
			})
		}

		transformedUsers = append(transformedUsers, map[string]interface{}{
			"id":             user.Id,
			"name":           user.Name,
			"username":       user.Username,
			"remember_token": user.RememberToken,
			"roles":          roles,
		})
	}

	// Response JSON dengan pagination info
	return ctx.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status: true,
		Data:   transformedUsers,
		Page:   page,
		Limit:  limit,
		Total:  totalRecords,
	})
}
