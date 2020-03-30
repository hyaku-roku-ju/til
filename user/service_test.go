package user

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type SpyDataSource struct {
	called        bool
	preferredTime PreferredTime
}

func (s *SpyDataSource) Create(ctx context.Context, preferredTime PreferredTime) (string, error) {
	s.called = true
	s.preferredTime = preferredTime
	return "1", nil
}

func (s *SpyDataSource) SetPreferredTime(ctx context.Context, id string, preferredTime PreferredTime) error {
	s.called = true
	s.preferredTime = preferredTime
	return nil
}

func TestCreateNewUser(t *testing.T) {
	var dataSource = &SpyDataSource{false, PreferredTime{0, 0}}

	ctx := context.Background()
	entity := UserEntity{dataSource}
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating new user"))
	}

	if dataSource.called != true {
		t.Errorf(fmt.Sprintf("CreateNewUser did not call data source"))
	}
}

func TestCreateNewUserPreferredTime(t *testing.T) {
	var dataSource = &SpyDataSource{false, PreferredTime{0, 0}}

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
	entity := UserEntity{dataSource}
	_, err := entity.CreateNewUser(ctx)

	if err != nil {
		t.Errorf(fmt.Sprintf("Error creating new user"))
	}

	preferredTime := PreferredTime{hour, min}

	if dataSource.preferredTime != preferredTime {
		t.Errorf(fmt.Sprintf("Error creating user with preferred time %d", dataSource.preferredTime))
	}
}
