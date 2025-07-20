package usecase

import (
	"context"
	"fmt"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/helpers"
	"time"
)

type authUsecase struct {
	userRepository        domain.UserRepository
	reffLookupRepository  domain.ReffLookupRepository
	whitelistIpRepository domain.WhitelistIpRepository
	contextTimeout        time.Duration
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepository domain.UserRepository, reffLookupRepository domain.ReffLookupRepository, whitelistIpRepository domain.WhitelistIpRepository, contextTimeout time.Duration) web.AuthUsecase {
	return &authUsecase{
		userRepository:        userRepository,
		reffLookupRepository:  reffLookupRepository,
		whitelistIpRepository: whitelistIpRepository,
		contextTimeout:        contextTimeout,
	}
}

// GetUserFromSession retrieves a user from a session ID
func (au *authUsecase) GetUserFromSession(ctx context.Context, sessionID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()

	return au.userRepository.GetSession(ctx, sessionID)
}

// DecryptSessionID decrypts a session cookie to get the session ID
func (au *authUsecase) DecryptSessionID(key, sessionCookie string) (string, error) {
	return helpers.DecryptAES(sessionCookie, key)
}

func (au *authUsecase) IpSetting(ctx context.Context, ip string) error {
	isIp, err := au.reffLookupRepository.GetDataByGroup(ctx, "ip_setting")
	if err != nil {
		panic(err)
	}

	if len(isIp) == 0 || isIp[0].LookupValue == "0" {
		return nil
	}

	ipData, err := au.whitelistIpRepository.GetDataByIp(ctx, ip)
	if err != nil || ipData.Ip != ip || !ipData.IsActive {
		return fmt.Errorf("ip %s not authorized", ip)
	}

	return nil
}
