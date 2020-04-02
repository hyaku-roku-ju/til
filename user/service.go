package user

import (
	"context"
	"time"
)

type UserEntity struct {
	dataSource UserDataSource
}

// Set time.Now outside of method for testing
var timeNow = time.Now

func NewEntity(dataSource UserDataSource) UserEntity {
  return UserEntity{dataSource}
}

func (u UserEntity) CreateNewUser(ctx context.Context) (string, error) {
	now := timeNow().UTC()
	preferredTime := PreferredTime{Hour: now.Hour(), Min: now.Minute()}
	id, err := u.dataSource.Create(ctx, preferredTime)
	return id, err
}

func (u UserEntity) SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error {
	err := u.dataSource.SetPreferredTime(ctx, id, preferredTime)
	return err
}
