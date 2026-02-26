package entity

import "gorm.io/gorm"

// User merepresentasikan entitas User di database.
// Struct ini digunakan oleh repository layer untuk operasi database.
type User struct {
	gorm.Model          // Embed gorm.Model (ID, CreatedAt, UpdatedAt, DeletedAt)
	Name       string   `json:"name" gorm:"not null"`
	Email      string   `json:"email" gorm:"uniqueIndex;not null"`
	Password   string   `json:"-" gorm:"not null"` // json:"-" agar tidak ter-serialize
	Age        int      `json:"age"`
}
