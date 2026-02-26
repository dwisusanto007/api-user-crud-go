package service_test

import (
	"api-user-crud-go/dto"
	"api-user-crud-go/entity"
	"api-user-crud-go/service"
	"errors"
	"testing"
)

// ==========================================
// MOCK REPOSITORY
// ==========================================

// mockUserRepo adalah implementasi mock dari repository.UserRepository.
type mockUserRepo struct {
	users  map[uint]*entity.User
	nextID uint
}

func newMockRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*entity.User), nextID: 1}
}

func (m *mockUserRepo) Create(user *entity.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindAll() ([]entity.User, error) {
	var result []entity.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockUserRepo) FindByID(id uint) (*entity.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) Update(user *entity.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(id uint) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

// ==========================================
// TESTS
// ==========================================

func newService() service.UserService {
	return service.NewUserService(newMockRepo())
}

func TestCreateUser(t *testing.T) {
	svc := newService()

	resp, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   25,
	})

	if err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}
	if resp.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if resp.Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", resp.Name)
	}
	if resp.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", resp.Email)
	}
	if resp.Age != 25 {
		t.Errorf("expected age 25, got %d", resp.Age)
	}
}

func TestGetAllUsers_Empty(t *testing.T) {
	svc := newService()

	users, err := svc.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers returned unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUsers_WithData(t *testing.T) {
	svc := newService()

	svc.CreateUser(dto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})
	svc.CreateUser(dto.CreateUserRequest{Name: "Bob", Email: "bob@example.com", Age: 30})

	users, err := svc.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers returned unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetUserByID_Found(t *testing.T) {
	svc := newService()

	created, _ := svc.CreateUser(dto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})

	user, err := svc.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned unexpected error: %v", err)
	}
	if user.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, user.ID)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.GetUserByID(999)
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	svc := newService()

	created, _ := svc.CreateUser(dto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})

	updated, err := svc.UpdateUser(created.ID, dto.UpdateUserRequest{
		Name: "Alice Updated",
		Age:  30,
	})
	if err != nil {
		t.Fatalf("UpdateUser returned unexpected error: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Errorf("expected name 'Alice Updated', got '%s'", updated.Name)
	}
	if updated.Age != 30 {
		t.Errorf("expected age 30, got %d", updated.Age)
	}
	// Email should remain unchanged
	if updated.Email != "alice@example.com" {
		t.Errorf("expected email unchanged, got '%s'", updated.Email)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.UpdateUser(999, dto.UpdateUserRequest{Name: "Ghost"})
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	svc := newService()

	created, _ := svc.CreateUser(dto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})

	err := svc.DeleteUser(created.ID)
	if err != nil {
		t.Fatalf("DeleteUser returned unexpected error: %v", err)
	}

	// Verify user is gone
	_, err = svc.GetUserByID(created.ID)
	if err == nil {
		t.Error("expected error after deletion, got nil")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc := newService()

	err := svc.DeleteUser(999)
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}
