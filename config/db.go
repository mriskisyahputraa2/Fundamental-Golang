package config

import (
	"Golang/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Function Untuk Terhubung ke Database
func ConnectDB() {
	// Membaca env dengan nama database_uri
	dsn := os.Getenv("DATABASE_URI")

	// Pengecekan jika ENV tidak ada maka akan fatal
	if dsn == "" {
		log.Fatal("Env Variabel Masih Belum Diisi")
	}

	// Membuat 2 variable untuk menampung koneksi dan error
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Pengecekan jika database tidak bisa terhubung
	if err != nil {
		log.Fatal("Gagal Terhubung ke Database", err)
	}

	// Melakukan Migrasi Database untuk membuat table
	err = database.AutoMigrate(&models.User{}, &models.Event{})

	// Pengecekan jika gagal Migrasi Database
	if err != nil {
		log.Fatal("Gagal Melakukan Migrasi Database", err)
	}

	DB = database
	log.Println("Berhasil Terkoneksi ke Database")
}
