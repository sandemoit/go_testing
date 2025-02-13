package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	amonia   uint = 1
	urea     uint = 2
	utilitas uint = 3
)

// @Summary Get Data Efisiensi Amoniak
// @Description Mengambil data Efisiensi Amoniak
// @Tags Efisiensi Dailies
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/getdata-effisiensi_amoniak [get]
func GetDataEfisiensiAmoniak(c *fiber.Ctx) error {
	// Ambil parameter page dan limit dari query string
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	items, totalItems, err := models.GetAllEfisiensi(database.DB, amonia, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to get data efisiensi amoniak: " + err.Error(),
		})
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to get data efisiensi amoniak",
		Data: map[string]interface{}{
			"items": items,
		},
		Limit:      limit,
		Page:       page,
		Total:      len(items),
		TotalPages: totalPages,
	})
}

// @Summary Get Data Efisiensi Urea
// @Description Mengambil data Efisiensi Urea
// @Tags Efisiensi Dailies
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/getdata-effisiensi_urea [get]
func GetDataEfisiensiUrea(c *fiber.Ctx) error {
	// Ambil parameter page dan limit dari query string
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	items, totalItems, err := models.GetAllEfisiensi(database.DB, urea, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to get data efisiensi amoniak: " + err.Error(),
		})
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to get data efisiensi amoniak",
		Data: map[string]interface{}{
			"items": items,
		},
		Limit:      limit,
		Page:       page,
		Total:      len(items),
		TotalPages: totalPages,
	})
}

// @Summary Get Data Efisiensi Utilitas
// @Description Mengambil data Efisiensi Utilitas
// @Tags Efisiensi Dailies
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/getdata-effisiensi_utilitas [get]
func GetDataEfisiensiUtilitas(c *fiber.Ctx) error {
	// Ambil parameter page dan limit dari query string
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	items, totalItems, err := models.GetAllEfisiensi(database.DB, utilitas, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Failed to get data efisiensi amoniak: " + err.Error(),
		})
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to get data efisiensi amoniak",
		Data: map[string]interface{}{
			"items": items,
		},
		Limit:      limit,
		Page:       page,
		Total:      len(items),
		TotalPages: totalPages,
	})
}
