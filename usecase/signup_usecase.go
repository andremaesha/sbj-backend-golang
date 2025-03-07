package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/curl"
	"sbj-backend/internal/helpers"
	"time"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{userRepository: userRepository, contextTimeout: contextTimeout}
}

func (su *signupUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.userRepository.Create(ctx, user)
}

func (su *signupUsecase) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
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

	url := ""

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
