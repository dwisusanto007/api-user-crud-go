package repository

import (
	"api-user-crud-go/entity"
	"errors"

	"gorm.io/gorm"
)

// UserRepository adalah interface untuk operasi database User.
// Menggunakan pattern repository untuk memisahkan logika data access.
type UserRepository interface {
	Create(user *entity.User) error
	FindAll() ([]entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}

// userRepositoryImpl adalah implementasi dari UserRepository.
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create menambahkan user baru ke database.
func (r *userRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindAll mengambil semua user dari database.
func (r *userRepositoryImpl) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

// FindByID mencari user berdasarkan ID.
func (r *userRepositoryImpl) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update mengupdate data user yang sudah ada.
func (r *userRepositoryImpl) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete menghapus user berdasarkan ID.
func (r *userRepositoryImpl) Delete(id uint) error {
	result := r.db.Delete(&entity.User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}
