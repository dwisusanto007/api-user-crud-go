package grpcserver

import (
	"api-user-crud-go/dto"
	"api-user-crud-go/proto"
	"api-user-crud-go/service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCServer mengimplementasikan UserServiceServer yang dihasilkan dari proto.
// Semua business logic didelegasikan ke UserService yang sudah ada.
type UserGRPCServer struct {
	proto.UnimplementedUserServiceServer
	userService service.UserService
}

// NewUserGRPCServer membuat instance baru UserGRPCServer.
func NewUserGRPCServer(userService service.UserService) *UserGRPCServer {
	return &UserGRPCServer{userService: userService}
}

// CreateUser menangani RPC CreateUser - membuat user baru.
func (s *UserGRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserMessage, error) {
	// Validasi input
	if req.Name == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}
	if req.Age <= 0 {
		return nil, status.Error(codes.InvalidArgument, "age must be greater than 0")
	}

	// Panggil service yang sudah ada
	resp, err := s.userService.CreateUser(dto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Age:   int(req.Age),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return toProtoUser(resp), nil
}

// GetAllUsers menangani RPC GetAllUsers - mengambil semua user.
func (s *UserGRPCServer) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	users, err := s.userService.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve users: %v", err)
	}

	var protoUsers []*proto.UserMessage
	for i := range users {
		protoUsers = append(protoUsers, toProtoUser(&users[i]))
	}

	return &proto.GetAllUsersResponse{Users: protoUsers}, nil
}

// GetUser menangani RPC GetUser - mengambil user berdasarkan ID.
func (s *UserGRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserMessage, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	user, err := s.userService.GetUserByID(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return toProtoUser(user), nil
}

// UpdateUser menangani RPC UpdateUser - mengupdate data user.
func (s *UserGRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserMessage, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	user, err := s.userService.UpdateUser(uint(req.Id), dto.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Age:   int(req.Age),
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to update user: %v", err)
	}

	return toProtoUser(user), nil
}

// DeleteUser menangani RPC DeleteUser - menghapus user berdasarkan ID.
func (s *UserGRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	err := s.userService.DeleteUser(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete user: %v", err)
	}

	return &proto.DeleteUserResponse{Message: "User deleted successfully"}, nil
}

// toProtoUser adalah helper untuk konversi dari dto.UserResponse ke proto.UserMessage.
func toProtoUser(u *dto.UserResponse) *proto.UserMessage {
	return &proto.UserMessage{
		Id:    uint32(u.ID),
		Name:  u.Name,
		Email: u.Email,
		Age:   int32(u.Age),
	}
}
