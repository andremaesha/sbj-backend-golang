package usecase

import (
	"context"
	"sbj-backend/domain"
	"time"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{userRepository: userRepository, contextTimeout: contextTimeout}
}

func (lu *loginUsecase) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUsecase) SetSession(c context.Context, expire int, idSession string, data *domain.User) error {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	lu.userRepository.SetExpire(expire)
	return lu.userRepository.SetSession(ctx, idSession, data)
}
