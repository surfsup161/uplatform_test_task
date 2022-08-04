package service

import (
	"context"
	"log"

	"github.com/surfsup161/uplatform_test_task/repo"

	"github.com/google/uuid"
	"github.com/surfsup161/uplatform_test_task/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUserService(repo repo.Repo) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	model.UserServiceServer
	repo repo.Repo
}

func (us *UserService) Set(_ context.Context, u *model.SetUserRequest) (*model.SetUserResponse, error) {
	log.Println("set user")
	userInfo := u.GetUser()
	if userInfo == nil {
		return nil, status.Errorf(codes.Aborted, "model info is empty")
	}

	userInfo.Id = uuid.New().String()
	err := us.repo.SetUser(userInfo)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
	}

	return &model.SetUserResponse{User: userInfo}, nil
}

func (us *UserService) Get(_ context.Context, u *model.GetUserRequest) (*model.GetUserResponse, error) {
	id := u.GetId()
	log.Println("get user", id)
	userInfo, err := us.repo.GetUser(id)
	if err != nil {
		return nil, err
	}

	if len(userInfo.GetId()) == 0 {
		return nil, status.Errorf(codes.NotFound, "model not found")
	}

	return &model.GetUserResponse{User: userInfo}, nil
}

func (us *UserService) Delete(_ context.Context, u *model.DeleteUserRequest) (*model.DeleteUserResponse, error) {
	id := u.GetId()
	log.Println("delete user", id)
	err := us.repo.DeleteUser(id)
	if err != nil {
		return nil, err
	}

	return &model.DeleteUserResponse{Id: id}, nil
}
