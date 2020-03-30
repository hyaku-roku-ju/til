package user

import (
	"fmt"
	"context"
	"time"
	"testing"
)

type SpyDataSource struct {
	called bool
	preferredTime int
}

func (s *SpyDataSource) Create(ctx context.Context, preferredTime int) (string, error) {
	s.called = true
	s.preferredTime = preferredTime
	return "1", nil
}

func (s *SpyDataSource) SetPreferredTime(ctx context.Context, id string, preferredTime int) error {
	s.called = true
	s.preferredTime = preferredTime
	return nil
}

func TestCreateNewUser(t *testing.T) {
	var dataSource = &SpyDataSource{ false, 0 }

	ctx := context.Background()
	entity := UserEntity{ dataSource }
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating new user"))
	}

	if dataSource.called != true {
		t.Errorf(fmt.Sprintf("CreateNewUser did not call data source"))
	}
}

func TestCreateNewUserPreferredTime(t *testing.T) {
	var dataSource = &SpyDataSource{ false, 0 }

	hour := 6
	min := 15
	timeNow = func() time.Time {
		t := time.Date(2019, 1, 1, hour, min, 0, 0, time.UTC)
		return t
	}

	defer func() {
		timeNow = time.Now
	}()

	ctx := context.Background()
	entity := UserEntity{ dataSource }
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating new user"))
	}

	preferredTime := hour * 60 + min

	if dataSource.preferredTime != preferredTime {
		t.Errorf(fmt.Sprintf("Error creating user with preferred time %d", dataSource.preferredTime))
	}
}