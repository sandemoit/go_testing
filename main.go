// @title           API Teknologi
// @version         1.0
// @description     API Kebutuhan untuk Portal Teknologi
// @host            localhost:3000
// @BasePath        /api
//
//go:debug x509negativeserial=1
package main

import (
	"go_teknologi/app/database"
	"go_teknologi/routes"
	Utils "go_teknologi/utils"
	"os"
)

func main() {
	// laod env
	Utils.LoadEnv()

	// Koneksi ke database
	database.ConnectDB()

	// Inisialisasi Fiber
	app := Utils.CreateApp()

	routes.SetupRoutes(app)

	// Jalankan server
	app.Listen(":" + os.Getenv("GOT_PORT"))
}
