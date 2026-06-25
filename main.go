package main

import (
	"log"
	"os"

	"github.com/AryaPramudya898/backend-sepatu-futsal.git/config"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/routes"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	config.InitFirebase()

	config.InitDatabase()

	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/v1/health", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
