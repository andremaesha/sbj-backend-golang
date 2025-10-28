package web

import (
	"context"
	"mime/multipart"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
)

type SignupRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Avatar    string `json:"avatar"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

type UploadAvatarResponse struct {
	Id        string `json:"id"`
	Filename  string `json:"filename"`
	UrlAvatar string `json:"url_avatar"`
	Message   string `json:"message"`
}

type SignupUsecase interface {
	Create(c context.Context, user *domain.User, avatarUrl string) error
	CreateUser(env *bootstrap.Env, c context.Context, request *SignupRequest) error
	GetUserByEmail(c context.Context, email string) (*domain.User, error)
	UploadAvatar(env *bootstrap.Env, fileHeader *multipart.FileHeader) (*domain.ResponseCloudinary, error)
}
