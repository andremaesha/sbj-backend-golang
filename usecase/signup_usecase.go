package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/curl"
	"sbj-backend/internal/encry"
	"sbj-backend/internal/helpers"
	"time"
)

type signupUsecase struct {
	userRepository   domain.UserRepository
	avatarRepository domain.AvatarRepository
	contextTimeout   time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, avatarRepository domain.AvatarRepository, contextTimeout time.Duration) web.SignupUsecase {
	return &signupUsecase{
		userRepository:   userRepository,
		avatarRepository: avatarRepository,
		contextTimeout:   contextTimeout,
	}
}

func (su *signupUsecase) Create(c context.Context, user *domain.User, avatarUrl string) error {
	now := time.Now()
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	err := su.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	avatar := &domain.Avatar{
		UserId: user.Id,
		Url:    avatarUrl,
	}
	err = su.avatarRepository.Create(ctx, avatar)
	if err != nil {
		return err
	}

	return su.userRepository.Update(ctx, &domain.User{
		Id:        user.Id,
		AvatarId:  avatar.Id,
		UpdatedAt: &now,
	})
}

func (su *signupUsecase) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
}

func (su *signupUsecase) CreateUser(c context.Context, request *web.SignupRequest) error {
	encryptedPassword, err := encry.HashPassword(request.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  encryptedPassword,
	}

	return su.Create(c, user, request.Avatar)
}

func (su *signupUsecase) UploadAvatar(env *bootstrap.Env, fileHeader *multipart.FileHeader) (*domain.ResponseCloudinary, error) {
	response := new(domain.ResponseCloudinary)

	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	folder := "avatar_sbj"
	publicId := uuid.New().String()

	formula := fmt.Sprintf("folder=%s&public_id=%s&timestamp=%s%s", folder, publicId, timeStamp, env.CloudinaryApiSecret)

	signature := helpers.GenerateSH1(formula)

	fileUpload, err := helpers.SaveTempFile(fileHeader, "uploads", "img_", "_backup")
	println(fileUpload)
	if err != nil {
		panic(err)
	}

	url := env.CloudinaryUrl

	cloudinary, err := curl.Curl[*domain.ResponseCloudinary]("POST", url, nil, nil, map[string]string{
		"file":      fileUpload,
		"api_key":   env.CloudinaryApiKey,
		"timestamp": timeStamp,
		"signature": signature,
		"public_id": publicId,
		"folder":    folder,
	}, response)
	if err != nil {
		panic(err)
	}

	return cloudinary, nil
}
