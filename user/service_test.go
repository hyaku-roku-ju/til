package user

import (
	"fmt"
	"context"
	"time"
	"testing"
)

type SpyStore struct {
	called bool
	preferredTime int
}

func (s *SpyStore) Create(ctx context.Context, preferredTime int) (string, error) {
	s.called = true
	s.preferredTime = preferredTime
	return "1", nil
}

func (s *SpyStore) SetPreferredTime(ctx context.Context, id string, preferredTime int) error {
	s.called = true
	s.preferredTime = preferredTime
	return nil
}

func TestCreateNewUser(t *testing.T) {
	var userStore = &SpyStore{ false, 0 }

	ctx := context.Background()
	entity := UserEntity{ userStore }
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating user"))
	}

	if userStore.called != true {
		t.Errorf(fmt.Sprintf("UserStore was not called with create"))
	}
}

func TestCreateNewUserPreferredTime(t *testing.T) {
	var userStore = &SpyStore{ false, 0 }

	hour := 6
	min := 15
	timeNow = func() time.Time {
		t := time.Date(2019, 1, 1, hour, min, 0, 0, time.UTC)
		return t
	}

	ctx := context.Background()
	entity := UserEntity{ userStore }
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating user"))
	}

	preferredTime := hour * 60 + min

	if userStore.preferredTime != preferredTime {
		t.Errorf(fmt.Sprintf("Error creating preferred time %d, not %d", userStore.preferredTime, preferredTime))
	}
}