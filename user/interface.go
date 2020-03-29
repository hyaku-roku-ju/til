package user

import (
	"context"
)

type UserDataSource interface {
	Create(ctx context.Context, preferredTime int) (string, error)
	SetPreferredTime(ctx context.Context, id string, preferredTime int) error
}