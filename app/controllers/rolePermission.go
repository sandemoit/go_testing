package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models" // sesuaikan dengan nama project Anda

	"github.com/gofiber/fiber/v2"
)

func CreateRolePermission(c *fiber.Ctx) error {
	// Parse request body
	var request models.RolePermission
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Mulai transaksi
	tx := database.DB.Begin()

	// Cek apakah role exists
	var role models.Role
	if err := tx.First(&role, request.RoleId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Role tidak ditemukan",
		})
	}

	// Cek apakah permission exists
	var permission models.Permission
	if err := tx.First(&permission, request.PermissionId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Permission tidak ditemukan",
		})
	}

	// Tambahkan permission ke role
	if err := tx.Model(&role).Association("Permissions").Append(&permission); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal menambahkan permission ke role",
		})
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to save data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to save data",
	})
}

func DeleteRolePermission(c *fiber.Ctx) error {
	id := c.Params("id")

	// Cek apakah role permission exists
	var rolePermission models.RolePermission
	if err := database.DB.First(&rolePermission, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "role permission not found",
		})
	}

	// Hapus role permission jika ditemukan
	if err := database.DB.Delete(&rolePermission).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to delete role permission: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to delete role permission",
	})
}

func GetRolePermission(c *fiber.Ctx) error {
	var rolePermissions []models.RolePermission

	// Load data dengan Preload
	result := database.DB.
		Preload("Permission.Roles").
		Preload("Role.Permissions").
		Preload("Role.Users.Roles").
		Find(&rolePermissions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal mendapatkan data role permission: " + result.Error.Error(),
		})
	}

	// **Transform data agar roles menjadi array string**
	var response []map[string]interface{}
	for _, rp := range rolePermissions {
		// Konversi `roles` ke []string
		var roleNames []string
		for _, role := range rp.Permission.Roles {
			roleNames = append(roleNames, role.Role)
		}

		// Konversi `user roles` ke []string
		var userList []map[string]interface{}
		for _, user := range rp.Role.Users {
			var userRoles []string
			for _, role := range user.Roles {
				userRoles = append(userRoles, role.Role)
			}

			userList = append(userList, map[string]interface{}{
				"id":             user.Id,
				"name":           user.Name,
				"remember_token": user.RememberToken,
				"roles":          userRoles,
				"username":       user.Username,
			})
		}

		// Tambahkan ke response JSON
		response = append(response, map[string]interface{}{
			"id": rp.Id,
			"permission": map[string]interface{}{
				"id":         rp.Permission.Id,
				"permission": rp.Permission.Permission,
				"roles":      roleNames, // Hanya array string
			},
			"permission_id": rp.PermissionId,
			"role": map[string]interface{}{
				"id":          rp.Role.Id,
				"role":        rp.Role.Role,
				"permissions": []string{"edit_perm"}, // Tambahkan jika ingin hanya array string
				"users":       userList,
			},
			"role_id": rp.RoleId,
		})
	}

	// Return hasil yang telah dikonversi
	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Berhasil mendapatkan data role permission",
		Data:    response,
	})
}

func GetRolePermissionById(c *fiber.Ctx) error {
	// Ambil role_id dari parameter
	roleID := c.Params("role_id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Parameter role_id tidak boleh kosong",
		})
	}

	var rolePermissions []models.RolePermission

	// Query dengan preload untuk mendapatkan data lengkap
	result := database.DB.
		Preload("Permission.Roles").
		Preload("Role.Permissions").
		Preload("Role.Users.Roles").
		Where("role_id = ?", roleID).
		Find(&rolePermissions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Gagal mendapatkan data role permission: " + result.Error.Error(),
		})
	}

	// **Transform data agar roles hanya berisi array string**
	var response []map[string]interface{}
	for _, rp := range rolePermissions {
		// Konversi `roles` di dalam permission ke array string
		var roleNames []string
		for _, role := range rp.Permission.Roles {
			roleNames = append(roleNames, role.Role)
		}

		// Konversi `user roles` ke array string
		var userList []map[string]interface{}
		for _, user := range rp.Role.Users {
			var userRoles []string
			for _, role := range user.Roles {
				userRoles = append(userRoles, role.Role)
			}

			userList = append(userList, map[string]interface{}{
				"id":             user.Id,
				"name":           user.Name,
				"remember_token": user.RememberToken,
				"roles":          userRoles,
				"username":       user.Username,
			})
		}

		// Tambahkan hasil transformasi ke response JSON
		response = append(response, map[string]interface{}{
			"id": rp.Id,
			"permission": map[string]interface{}{
				"id":         rp.Permission.Id,
				"permission": rp.Permission.Permission,
				"roles":      roleNames, // Hanya array string
			},
			"permission_id": rp.PermissionId,
			"role": map[string]interface{}{
				"id":          rp.Role.Id,
				"role":        rp.Role.Role,
				"permissions": []string{"edit_perm"}, // Jika ingin hanya array string
				"users":       userList,
			},
			"role_id": rp.RoleId,
		})
	}

	// Return hasil JSON yang sudah diperbaiki
	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Berhasil mendapatkan data role permission",
		Data:    response,
	})
}
