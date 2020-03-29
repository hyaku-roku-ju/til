package user

import (
	"context"
	"time"
)

type UserEntity struct {
	dataSource UserDataSource
}

// Create timeNow for testing
var timeNow = time.Now

func (u UserEntity) CreateNewUser(ctx context.Context) (string, error) {
	preferredTime := timeNow().UTC().Hour() * 60 + timeNow().UTC().Minute()
	id, err := u.dataSource.Create(ctx, preferredTime)
	return id, err
}
