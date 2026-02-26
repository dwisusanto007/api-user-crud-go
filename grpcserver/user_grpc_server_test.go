package grpcserver_test

import (
	"api-user-crud-go/dto"
	"api-user-crud-go/entity"
	grpcserver "api-user-crud-go/grpcserver"
	"api-user-crud-go/proto"
	"api-user-crud-go/service"
	"context"
	"errors"
	"testing"
)

// ==========================================
// MOCK REPOSITORY (in-memory)
// ==========================================

type mockRepo struct {
	users  map[uint]*entity.User
	nextID uint
}

func newMockRepo() *mockRepo {
	return &mockRepo{users: make(map[uint]*entity.User), nextID: 1}
}

func (m *mockRepo) Create(user *entity.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *mockRepo) FindAll() ([]entity.User, error) {
	var result []entity.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockRepo) FindByID(id uint) (*entity.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockRepo) Update(user *entity.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepo) Delete(id uint) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

// newServer membuat gRPC server baru dengan mock repo untuk setiap test.
func newServer() *grpcserver.UserGRPCServer {
	svc := service.NewUserService(newMockRepo())
	return grpcserver.NewUserGRPCServer(svc)
}

var ctx = context.Background()

// ==========================================
// TESTS: CreateUser
// ==========================================

func TestGRPC_CreateUser_Success(t *testing.T) {
	srv := newServer()

	resp, err := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   25,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Id == 0 {
		t.Error("expected non-zero ID")
	}
	if resp.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", resp.Name)
	}
	if resp.Age != 25 {
		t.Errorf("expected age 25, got %d", resp.Age)
	}
}

func TestGRPC_CreateUser_MissingName(t *testing.T) {
	srv := newServer()

	_, err := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Email: "alice@example.com",
		Age:   25,
	})

	if err == nil {
		t.Error("expected validation error for missing name, got nil")
	}
}

func TestGRPC_CreateUser_MissingEmail(t *testing.T) {
	srv := newServer()

	_, err := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name: "Alice",
		Age:  25,
	})

	if err == nil {
		t.Error("expected validation error for missing email, got nil")
	}
}

func TestGRPC_CreateUser_InvalidAge(t *testing.T) {
	srv := newServer()

	_, err := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   0,
	})

	if err == nil {
		t.Error("expected validation error for age=0, got nil")
	}
}

// ==========================================
// TESTS: GetAllUsers
// ==========================================

func TestGRPC_GetAllUsers_Empty(t *testing.T) {
	srv := newServer()

	resp, err := srv.GetAllUsers(ctx, &proto.GetAllUsersRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Users) != 0 {
		t.Errorf("expected 0 users, got %d", len(resp.Users))
	}
}

func TestGRPC_GetAllUsers_WithData(t *testing.T) {
	srv := newServer()

	srv.CreateUser(ctx, &proto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})
	srv.CreateUser(ctx, &proto.CreateUserRequest{Name: "Bob", Email: "bob@example.com", Age: 30})

	resp, err := srv.GetAllUsers(ctx, &proto.GetAllUsersRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(resp.Users))
	}
}

// ==========================================
// TESTS: GetUser
// ==========================================

func TestGRPC_GetUser_Found(t *testing.T) {
	srv := newServer()

	created, _ := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name: "Alice", Email: "alice@example.com", Age: 25,
	})

	resp, err := srv.GetUser(ctx, &proto.GetUserRequest{Id: created.Id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Id != created.Id {
		t.Errorf("expected ID %d, got %d", created.Id, resp.Id)
	}
}

func TestGRPC_GetUser_NotFound(t *testing.T) {
	srv := newServer()

	_, err := srv.GetUser(ctx, &proto.GetUserRequest{Id: 999})
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestGRPC_GetUser_InvalidID(t *testing.T) {
	srv := newServer()

	_, err := srv.GetUser(ctx, &proto.GetUserRequest{Id: 0})
	if err == nil {
		t.Error("expected validation error for id=0, got nil")
	}
}

// ==========================================
// TESTS: UpdateUser
// ==========================================

func TestGRPC_UpdateUser_Success(t *testing.T) {
	srv := newServer()

	created, _ := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name: "Alice", Email: "alice@example.com", Age: 25,
	})

	resp, err := srv.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:   created.Id,
		Name: "Alice Updated",
		Age:  30,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "Alice Updated" {
		t.Errorf("expected 'Alice Updated', got '%s'", resp.Name)
	}
	if resp.Age != 30 {
		t.Errorf("expected age 30, got %d", resp.Age)
	}
}

func TestGRPC_UpdateUser_NotFound(t *testing.T) {
	srv := newServer()

	_, err := srv.UpdateUser(ctx, &proto.UpdateUserRequest{Id: 999, Name: "Ghost"})
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

// ==========================================
// TESTS: DeleteUser
// ==========================================

func TestGRPC_DeleteUser_Success(t *testing.T) {
	srv := newServer()

	created, _ := srv.CreateUser(ctx, &proto.CreateUserRequest{
		Name: "Alice", Email: "alice@example.com", Age: 25,
	})

	resp, err := srv.DeleteUser(ctx, &proto.DeleteUserRequest{Id: created.Id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}

	// Verify user is gone
	_, err = srv.GetUser(ctx, &proto.GetUserRequest{Id: created.Id})
	if err == nil {
		t.Error("expected error after deletion, got nil")
	}
}

func TestGRPC_DeleteUser_NotFound(t *testing.T) {
	srv := newServer()

	_, err := srv.DeleteUser(ctx, &proto.DeleteUserRequest{Id: 999})
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestGRPC_DeleteUser_InvalidID(t *testing.T) {
	srv := newServer()

	_, err := srv.DeleteUser(ctx, &proto.DeleteUserRequest{Id: 0})
	if err == nil {
		t.Error("expected validation error for id=0, got nil")
	}
}

// ==========================================
// HELPER: ensure UserResponse implements dto
// ==========================================

var _ *dto.UserResponse = (*dto.UserResponse)(nil)
