package learning

import (
	"context"
)

type LearningRepository interface {
  StoreLearning(ctx context.Context, learning Learning) (string, error)
}

type LearningService struct {
  learningRepository LearningRepository 
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
