package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sbj-backend/domain"
	"strconv"
	"time"
)

type userRepository struct {
	db          *gorm.DB
	redis       *redis.Client
	table       string
	redisPrefix []string
	expire      int
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client, table string, redisPrefix ...string) domain.UserRepository {
	return &userRepository{db: db, redis: redisClient, table: table, redisPrefix: redisPrefix}
}

func (ur *userRepository) SetExpire(expire int) {
	ur.expire = expire
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	return ur.db.WithContext(ctx).Table(ur.table).Create(user).Error
}

func (ur *userRepository) Update(ctx context.Context, user *domain.User) error {
	return ur.db.WithContext(ctx).Table(ur.table).Updates(user).Error
}

func (ur *userRepository) GetSession(ctx context.Context, idSession string) (*domain.User, error) {
	prefix := ur.redisPrefix[0]
	key := prefix + idSession

	data, err := ur.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Check if data is empty
	if len(data) == 0 {
		return nil, fmt.Errorf("session not found")
	}

	var user domain.User

	// Map Redis hash fields to User struct
	if id, exists := data["user_id"]; exists {
		if parsedId, err2 := strconv.Atoi(id); err2 == nil {
			user.Id = parsedId
		}
	}

	if email, exists := data["email"]; exists {
		user.Email = email
	}

	if role, exists := data["role"]; exists {
		user.Role = role
	}

	return &user, nil
}

func (ur *userRepository) SetSession(ctx context.Context, idSession string, user *domain.User) error {
	prefix := ur.redisPrefix[0]
	key := prefix + idSession
	timeExpire := time.Minute * time.Duration(ur.expire)

	err := ur.redis.HSet(ctx, key, user).Err()
	if err != nil {
		return err
	}

	return ur.redis.Expire(ctx, key, timeExpire).Err()
}

func (ur *userRepository) DeleteSession(ctx context.Context, idSession string) (int64, error) {
	prefix := ur.redisPrefix[0]
	key := prefix + idSession

	return ur.redis.Del(ctx, key).Result()
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (*domain.User, error) {
	result := new(domain.User)

	err := ur.db.WithContext(c).Table(ur.table).Where("email = ?", email).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
