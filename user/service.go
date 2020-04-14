package user

import (
	"context"
	"time"
)

type UserEntity struct {
	repository UserRepository
}

// Set time.Now outside of method for testing
var timeNow = time.Now

func NewEntity(repository UserRepository) UserEntity {
  return UserEntity{repository}
}

func (u UserEntity) CreateNewUser(ctx context.Context) (string, error) {
	now := timeNow().UTC()
	preferredTime := PreferredTime{Hour: now.Hour(), Min: now.Minute()}
	id, err := u.repository.Create(ctx, preferredTime)
	return id, err
}

func (u UserEntity) SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error {
	err := u.repository.SetPreferredTime(ctx, id, preferredTime)
	return err
}
