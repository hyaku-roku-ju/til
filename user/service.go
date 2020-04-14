package user

import (
	"context"
	"time"
)

type UserEntity struct {
	userRepository UserRepository
}

// Set time.Now outside of method for testing
var timeNow = time.Now

func NewEntity(userRepository UserRepository) UserEntity {
  return UserEntity{userRepository}
}

func (u UserEntity) CreateNewUser(ctx context.Context) (string, error) {
	now := timeNow().UTC()
	preferredTime := PreferredTime{Hour: now.Hour(), Min: now.Minute()}
	id, err := u.userRepository.Create(ctx, preferredTime)
	return id, err
}

func (u UserEntity) SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error {
	err := u.userRepository.SetPreferredTime(ctx, id, preferredTime)
	return err
}
