package learning

import (
	"context"
	"math/rand"
	"time"
)

type LearningRepository interface {
	StoreLearning(ctx context.Context, learning Learning) (string, error)
	CountConfirmedLearnings(ctx context.Context, reporterId string) (int, error)
	GetConfirmedLearning(ctx context.Context, reporterId string, skip int) (Learning, error)
}

type LearningService struct {
	learningRepository LearningRepository
}

func NewService(learningRepository LearningRepository) LearningService {
	return LearningService{learningRepository}
}

func (self *LearningService) GetRandomLearning(ctx context.Context, reporterId string) (Learning, error) {
	success := make(chan Learning)
	fail := make(chan error)

	go func() {
		count, err := self.learningRepository.CountConfirmedLearnings(ctx, reporterId)
		if err != nil {
			fail <- err
		}

		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		learningsToSkip := rnd.Intn(count)
		learning, err := self.learningRepository.GetConfirmedLearning(ctx, reporterId, learningsToSkip)
		if err != nil {
			fail <- err
		}
		success <- learning
	}()

	select {
	case <-ctx.Done():
		return Learning{}, ctx.Err()
	case err := <-fail:
		return Learning{}, err
	case learning := <-success:
		return learning, nil
	}
}

func (self *LearningService) StoreLearning(ctx context.Context, learning Learning) (learningId string, err error) {
	if _, err := learning.IsValid(); err != nil {
		return "", err
	}

	fail := make(chan error)
	success := make(chan string)

	go func() {
		id, err := self.learningRepository.StoreLearning(ctx, learning)
		if err != nil {
			fail <- err
		}
		success <- id
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-fail:
		return "", err
	case learningId := <-success:
		return learningId, nil
	}
}
