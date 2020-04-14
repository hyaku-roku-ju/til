package learning

import (
	"context"
)

type LearningRepository interface {
  StoreLearning(ctx context.Context, learning Learning) (string, error)
  GetRandomLearning(ctx context.Context, reporterId string) (Learning, error)
}

type LearningService struct {
  learningRepository LearningRepository
}

func NewService(repo LearningRepository) LearningService {
  return LearningService{repo}
}

func (self *LearningService) GetRandomLearning(ctx context.Context, reporterId string) (Learning, error) {
  complete := make(chan Learning)
  fail := make(chan error)

  go func() {
    learning, err := self.learningRepository.GetRandomLearning(ctx, reporterId)
    if err != nil {
      fail<-err
    }
    complete<-learning
  }()

  select {
    case <-ctx.Done():
      return Learning{}, ctx.Err()
    case err := <-fail:
      return Learning{}, err
    case learning := <-complete:
      return learning, nil
  }
}

func (self *LearningService) StoreLearning(ctx context.Context, learning Learning) (learningId string, err error) {
  if _, err := learning.IsValid(); err != nil {
    return "", err
  }

  fail := make(chan error)
  complete := make(chan string)

  go func() {
    id, err := self.learningRepository.StoreLearning(ctx, learning)
    if err != nil {
      fail<-err
    }
    complete<-id
  }()

  select {
    case <-ctx.Done():
      return "", ctx.Err()
    case err := <-fail:
      return "", err
    case learningId := <-complete:
      return learningId, nil
  }
}
