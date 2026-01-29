package service

import (
	"api-user-crud-go/dto"
	"api-user-crud-go/entity"
	"api-user-crud-go/repository"
)

// UserService adalah interface untuk business logic User.
// Layer ini menangani konversi antara DTO dan Entity.
type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
}

// userServiceImpl adalah implementasi dari UserService.
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService membuat instance baru UserService.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

// CreateUser menambahkan user baru.
func (s *userServiceImpl) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Konversi dari DTO ke Entity
	user := &entity.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	// Simpan ke database melalui repository
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Konversi dari Entity ke DTO Response
	return toUserResponse(user), nil
}

// GetAllUsers mengambil semua user.
func (s *userServiceImpl) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Konversi slice Entity ke slice DTO
	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *toUserResponse(&user))
	}

	return responses, nil
}

// GetUserByID mengambil user berdasarkan ID.
func (s *userServiceImpl) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// UpdateUser mengupdate data user.
func (s *userServiceImpl) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Cek apakah user ada
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update field yang diisi (non-empty)
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	// Simpan perubahan
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// DeleteUser menghapus user berdasarkan ID.
func (s *userServiceImpl) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// toUserResponse adalah helper function untuk konversi Entity ke DTO Response.
func toUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
