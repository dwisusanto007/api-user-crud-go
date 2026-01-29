package config

import (
	"api-user-crud-go/entity"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB menginisialisasi koneksi database dan melakukan migrasi.
// Fungsi ini mengembalikan instance *gorm.DB untuk digunakan di layer lain.
func InitDB() *gorm.DB {
	// Membuka koneksi ke SQLite database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Auto-migrate schema
	// GORM akan membuat/update tabel berdasarkan struct entity
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	log.Println("✓ Database terkoneksi & migrasi berhasil!")
	return db
}
