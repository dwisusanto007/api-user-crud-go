package service

import (
	"api-user-crud-go/config"
	"api-user-crud-go/dto"
	"api-user-crud-go/entity"
	"api-user-crud-go/middleware"
	"api-user-crud-go/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService adalah interface untuk authentication logic
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

// authServiceImpl adalah implementasi dari AuthService
type authServiceImpl struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register mendaftarkan user baru
func (s *authServiceImpl) Register(req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Cek apakah email sudah terdaftar
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Buat user baru
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Age:      req.Age,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		},
	}, nil
}
func (s *authServiceImpl) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		},
	}, nil
}
