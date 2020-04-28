package user

import (
	"context"
	"time"
)

type UserService struct {
	userRepository UserRepository
}

// Set time.Now outside of method for testing
var timeNow = time.Now

func NewService(userRepository UserRepository) UserService {
  return UserService{userRepository}
}

func (u UserService) CreateNewUser(ctx context.Context) (string, error) {
	now := timeNow().UTC()
	preferredTime := PreferredTime{Hour: now.Hour(), Min: now.Minute()}
	id, err := u.userRepository.Create(ctx, preferredTime)
	return id, err
}

func (u UserService) SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error {
	err := u.userRepository.SetPreferredTime(ctx, id, preferredTime)
	return err
}
