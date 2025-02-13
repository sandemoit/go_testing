package controllers

import (
	"go_teknologi/app/database"
	"go_teknologi/app/models"
	"go_teknologi/app/repositories"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// Worker Pool Configuration
const workerCount = 5 // Jumlah worker untuk memproses data secara paralel

// @Summary Import Excel Data
// @Description Mengimpor data dari file Excel dan menyimpannya ke database
// @Tags Data Produksi
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File Excel yang akan diimpor"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 400 {object} models.GlobalErrorHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/import-prod_dailies [post]
func ImportExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}
	defer src.Close()

	excel, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}

	sheetName := excel.GetSheetName(0)
	rows, err := excel.GetRows(sheetName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Channel untuk mengirim data ke worker
	dataChan := make(chan []string, len(rows)-1)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var records []models.ProdDailies

	// Worker pool untuk pemrosesan data
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go processExcelRow(&wg, dataChan, &records, &mu)
	}

	// Kirim data ke channel untuk diproses oleh worker
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		dataChan <- row
	}
	close(dataChan) // Tutup channel setelah semua data dikirim

	// Tunggu semua worker selesai
	wg.Wait()

	// Masukkan data ke database secara batch
	if len(records) > 0 {
		database.DB.Create(&records)
	}

	return c.JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Data imported successfully",
	})
}

// @Summary Get Draft Data
// @Description Mengambil data produksi dengan status draft
// @Tags Data Produksi
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param tanggal_prod query string false "Tanggal periode (format: YYYY-MM-DD)"
// @Param tanggal_awal query string false "Tanggal awal periode (format: YYYY-MM-DD)"
// @Param tanggal_akhir query string false "Tanggal akhir periode (format: YYYY-MM-DD)"
// @Param kode_material query string false "Kode Material"
// @Param status query int false "Status"
// @Success 200 {object} models.GlobalSuccessHandlerResp{data=[]models.ProdDailies}
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/getdraft-prod_dailies [get]
func GetDataDraft(c *fiber.Ctx) error {
	return getData(c, 0)
}

// @Summary Get Released Data
// @Description Mengambil data produksi dengan status rilis
// @Tags Data Produksi
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param tanggal_prod query string false "Tanggal periode (format: YYYY-MM-DD)"
// @Param tanggal_awal query string false "Tanggal awal periode (format: YYYY-MM-DD)"
// @Param tanggal_akhir query string false "Tanggal akhir periode (format: YYYY-MM-DD)"
// @Param kode_material query string false "Kode Material"
// @Param status query int false "Status"
// @Success 200 {object} models.GlobalSuccessHandlerResp{data=[]models.ProdDailies}
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/getrilis-prod_dailies [get]
func GetDataRilis(c *fiber.Ctx) error {
	return getData(c, 1)
}

// @Summary Update Production Data
// @Description Mengupdate data produksi berdasarkan metode tertentu
// @Tags Data Produksi
// @Accept json
// @Produce json
// @Param request body models.UpdateProdDailiesRequest true "Payload untuk update data"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 400 {object} models.GlobalErrorHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/update-prod_dailies [put]
func UpdateProd(c *fiber.Ctx) error {
	var payload models.UpdateProdDailiesRequest

	// Parsing JSON request
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid request payload",
		})
	}

	// Validasi jika data kosong
	if len(payload.Data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Data cannot be empty",
		})
	}

	// Optimasi dengan map agar tidak ada UUID duplikat
	uuidSet := make(map[string]struct{})
	statusSet := make(map[uint]struct{})

	for _, item := range payload.Data {
		uuidSet[item.UUID] = struct{}{}
		statusSet[item.Status] = struct{}{}
	}

	// Convert map ke slice untuk query database
	var uuids []string
	for uuid := range uuidSet {
		uuids = append(uuids, uuid)
	}

	// Proses berdasarkan method
	switch payload.Method {
	case "delete":
		// Batch update berdasarkan UUID (soft delete)
		if err := database.DB.Model(&models.ProdDailies{}).
			Where("uuid IN (?)", uuids).
			Update("deleted_at", time.Now()).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
				Status:  false,
				Message: "Failed to delete prod dailies: " + err.Error(),
			})
		}

	case "rilis":
		// Pastikan hanya status yang valid (0 atau 1)
		if len(statusSet) > 1 || (len(statusSet) == 1 && !containsValidStatus(statusSet)) {
			return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
				Status:  false,
				Message: "Invalid status value",
			})
		}

		// Batch update status ke "1" untuk semua data yang sesuai
		if err := database.DB.Model(&models.ProdDailies{}).
			Where("uuid IN (?)", uuids).
			Update("status", 1).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
				Status:  false,
				Message: "Failed to update status prod dailies: " + err.Error(),
			})
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid method",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to update prod dailies",
	})
}

// @Summary Revoke Production Data
// @Description Mengupdate status produksi berdasarkan metode tertentu
// @Tags Data Produksi
// @Accept json
// @Produce json
// @Param request body models.UpdateProdDailiesRequest true "Payload untuk update status"
// @Success 200 {object} models.GlobalSuccessHandlerResp
// @Failure 400 {object} models.GlobalErrorHandlerResp
// @Failure 500 {object} models.GlobalErrorHandlerResp
// @Router /api/revoke-prod_dailies [put]
func RevokeProd(c *fiber.Ctx) error {
	var payload models.UpdateProdDailiesRequest

	// Parsing JSON request
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid request payload",
		})
	}

	// Validasi jika data kosong
	if len(payload.Data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Data cannot be empty",
		})
	}

	// Optimasi dengan map agar tidak ada UUID duplikat
	uuidSet := make(map[string]struct{})
	statusSet := make(map[uint]struct{})

	for _, item := range payload.Data {
		uuidSet[item.UUID] = struct{}{}
		statusSet[item.Status] = struct{}{}
	}

	// Convert map ke slice untuk query database
	var uuids []string
	for uuid := range uuidSet {
		uuids = append(uuids, uuid)
	}

	// Proses berdasarkan method
	switch payload.Method {
	case "revoke":
		// Batch update berdasarkan UUID (soft delete)
		if err := database.DB.Model(&models.ProdDailies{}).
			Where("uuid IN (?)", uuids).
			Update("status", 0).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
				Status:  false,
				Message: "Failed to delete prod dailies: " + err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: "Invalid method",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to revoke prod dailies",
	})
}

func containsValidStatus(statusSet map[uint]struct{}) bool {
	for status := range statusSet {
		if status != 0 && status != 1 {
			return false
		}
	}
	return true
}

func processExcelRow(wg *sync.WaitGroup, dataChan <-chan []string, records *[]models.ProdDailies, mu *sync.Mutex) {
	defer wg.Done()

	for row := range dataChan {
		jmlProd, _ := strconv.Atoi(row[4])

		// Validate Material
		var mat models.Material
		if err := database.DB.First(&mat, "kode_mat = ? AND material = ?", row[1], row[2]).Error; err != nil {
			continue
		}

		// Validate Plant
		var plant models.Plant
		if err := database.DB.First(&plant, "plant = ?", row[3]).Error; err != nil {
			continue
		}

		prodData := models.ProdDailies{
			TglProd:      row[0],
			UUID:         uuid.New().String(),
			KodeMaterial: row[1],
			WCID:         row[3],
			Status:       0,
			CreatedBy:    "060492",
			JmlProd:      float64(jmlProd),
		}

		// Prevent race condition when appending to slice
		mu.Lock()
		*records = append(*records, prodData)
		mu.Unlock()
	}
}

func getData(c *fiber.Ctx, status int) error {
	// Get query parameters from URL
	limit, _ := strconv.Atoi(c.Query("limit", "25"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	tglProd := c.Query("tgl_prod")
	kodeMaterial := c.Query("kode_material")
	wcID := c.Query("wc_id")
	tglAwal := c.Query("tanggal_awal")
	tglAkhir := c.Query("tanggal_akhir")

	// Validate limit & page
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	// Build base query with relationships
	query := repositories.QueryScope(status)

	// Apply filters if parameters exist
	if tglProd != "" {
		_, err := time.Parse("2006-01-02", tglProd)
		if err == nil {
			query = query.Where("prod_dailies.tgl_prod = ?", tglProd)
		}
	}
	if kodeMaterial != "" {
		query = query.Where("prod_dailies.kode_material = ?", kodeMaterial)
	}
	if wcID != "" {
		query = query.Where("prod_dailies.wc_id = ?", wcID)
	}
	if tglAwal != "" && tglAkhir != "" {
		startDate, errStart := time.Parse("2006-01-02", tglAwal)
		endDate, errEnd := time.Parse("2006-01-02", tglAkhir)
		if errStart == nil && errEnd == nil {
			query = query.Where("prod_dailies.tgl_prod BETWEEN ? AND ?", startDate, endDate)
		}
	}

	// Execute query with pagination
	var prodDailies []models.ProdDailies
	result := query.Limit(limit).Offset(offset).Find(&prodDailies)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.GlobalErrorHandlerResp{
			Status:  false,
			Message: result.Error.Error(),
		})
	}

	// Count total records for pagination
	var total int64
	query.Count(&total)

	return c.JSON(models.GlobalSuccessHandlerResp{
		Status:  true,
		Message: "Success to get prod dailies",
		Data:    prodDailies,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}
