package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"

	"github.com/gofiber/fiber/v2"
)

func CreateUserRole(ctx *fiber.Ctx) error {
	// Parse request body
	var request models.UserRole
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}
	// Mulai transaksi
	tx := database.DB.Begin()

	// Cek apakah role exists
	var role models.Role
	if err := tx.First(&role, request.RoleID).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Role tidak ditemukan",
		})
	}

	// Cek apakah user exists
	var user models.User
	if err := tx.First(&user, request.UserID).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "User tidak ditemukan",
		})
	}

	// Tambahkan role ke user
	if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal menambahkan role ke user",
		})
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to save data: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to save data",
	})
}

func DeleteUserRole(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Cek apakah user role exists
	var userRole models.UserRole
	if err := database.DB.First(&userRole, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "User role not found",
		})
	}

	// Hapus user role jika ditemukan
	if err := database.DB.Delete(&userRole).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to delete user role: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to delete user role",
	})
}

func GetUserRole(ctx *fiber.Ctx) error {
	var userRoles []models.UserRole

	// Query dengan preload untuk mendapatkan data sesuai dengan format JSON yang diharapkan
	result := database.DB.
		Preload("Role.Permissions.Roles").
		Preload("Role.Users.Roles").
		Preload("User.Roles").
		Find(&userRoles)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal mendapatkan data user role: " + result.Error.Error(),
		})
	}

	var response []map[string]interface{}

	for _, ur := range userRoles {
		// Konversi role permissions ke format yang diinginkan
		var permissions []map[string]interface{}
		for _, perm := range ur.Role.Permissions {
			var rolesList []string
			for _, r := range perm.Roles {
				rolesList = append(rolesList, r.Role)
			}

			permissions = append(permissions, map[string]interface{}{
				"id":         perm.Id,
				"permission": perm.Permission,
				"roles":      rolesList,
			})
		}

		// Konversi user roles ke []string
		var userRoles []string
		for _, role := range ur.User.Roles {
			userRoles = append(userRoles, role.Role)
		}

		// Konversi daftar users di dalam role
		var userList []map[string]interface{}
		for _, user := range ur.Role.Users {
			var userRoleList []string
			for _, r := range user.Roles {
				userRoleList = append(userRoleList, r.Role)
			}

			userList = append(userList, map[string]interface{}{
				"id":             user.Id,
				"name":           user.Name,
				"remember_token": user.RememberToken,
				"roles":          userRoleList,
				"username":       user.Username,
			})
		}

		// Memperbarui response JSON dengan format permissions yang baru
		response = append(response, map[string]interface{}{
			"id": ur.Id,
			"role": map[string]interface{}{
				"id":          ur.Role.Id,
				"permissions": permissions,
				"role":        ur.Role.Role,
				"users":       userList,
			},
			"role_id": ur.RoleID,
			"user": map[string]interface{}{
				"id":             ur.User.Id,
				"name":           ur.User.Name,
				"remember_token": ur.User.RememberToken,
				"roles":          userRoles,
				"username":       ur.User.Username,
			},
			"user_id": ur.UserID,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Berhasil mendapatkan data user role",
		Data:    response,
	})
}

func GetUserRoleById(ctx *fiber.Ctx) error {
	// Mendapatkan user_id dari parameter
	userID := ctx.Params("user_id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Parameter user_id tidak boleh kosong",
		})
	}

	var userRole models.UserRole

	// Query dengan preload untuk mendapatkan data role dan user berdasarkan user_id
	result := database.DB.
		Preload("Role.Permissions.Roles").
		Preload("Role.Users.Roles").
		Preload("User.Roles").
		Where("user_id = ?", userID).
		First(&userRole)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Data user role tidak ditemukan",
		})
	}

	// Konversi permissions ke format yang diinginkan
	var permissions []map[string]interface{}
	for _, perm := range userRole.Role.Permissions {
		var rolesList []string
		for _, r := range perm.Roles {
			rolesList = append(rolesList, r.Role)
		}

		permissions = append(permissions, map[string]interface{}{
			"id":         perm.Id,
			"permission": perm.Permission,
			"roles":      rolesList,
		})
	}

	// Konversi user roles ke []string
	var userRoles []string
	for _, role := range userRole.User.Roles {
		userRoles = append(userRoles, role.Role)
	}

	// Konversi daftar `users` di dalam role
	var userList []map[string]interface{}
	for _, user := range userRole.Role.Users {
		var userRoleList []string
		for _, r := range user.Roles {
			userRoleList = append(userRoleList, r.Role)
		}

		userList = append(userList, map[string]interface{}{
			"id":             user.Id,
			"name":           user.Name,
			"remember_token": user.RememberToken,
			"roles":          userRoleList,
			"username":       user.Username,
		})
	}

	// Membangun response JSON sesuai format yang diharapkan
	response := map[string]interface{}{
		"id": userRole.Id,
		"role": map[string]interface{}{
			"id":          userRole.Role.Id,
			"permissions": permissions,
			"role":        userRole.Role.Role,
			"users":       userList,
		},
		"role_id": userRole.RoleID,
		"user": map[string]interface{}{
			"id":             userRole.User.Id,
			"name":           userRole.User.Name,
			"remember_token": userRole.User.RememberToken,
			"roles":          userRoles,
			"username":       userRole.User.Username,
		},
		"user_id": userRole.UserID,
	}

	return ctx.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Berhasil mendapatkan data user role",
		Data:    response,
	})
}
