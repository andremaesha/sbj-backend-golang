package usecase

import (
	"context"
	"errors"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/helpers"
	"time"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) web.LoginUsecase {
	return &loginUsecase{userRepository: userRepository, contextTimeout: contextTimeout}
}

func (lu *loginUsecase) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUsecase) ValidateUserVerified(verified bool) error {
	if verified {
		return nil
	}

	return errors.New("user not verified, please verify your email first")
}

func (lu *loginUsecase) SetSession(c context.Context, expire int, idSession string, data *domain.User) error {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	lu.userRepository.SetExpire(expire)
	return lu.userRepository.SetSession(ctx, idSession, data)
}

func (lu *loginUsecase) EncryptSession(key, data string) string {
	content, err := helpers.EncryptAES(data, key)
	if err != nil {
		panic(err)
	}

	return content
}
