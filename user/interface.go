package user

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, preferredTime PreferredTime) (string, error)
	SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error
}
