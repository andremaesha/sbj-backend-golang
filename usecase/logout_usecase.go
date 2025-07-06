package usecase

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/helpers"
	"time"
)

type LogoutUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLogoutUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) web.LogoutUsecase {
	return &LogoutUsecase{userRepository: userRepository, contextTimeout: contextTimeout}
}

func (lu *LogoutUsecase) DeleteSession(c context.Context, idSession string) error {
	_, err := lu.userRepository.DeleteSession(c, idSession)
	return err
}

func (lu *LogoutUsecase) DecryptSession(key, data string) string {
	content, err := helpers.DecryptAES(data, key)
	if err != nil {
		panic(err)
	}

	return content
}
