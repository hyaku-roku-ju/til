package learning

import (
	"context"
	"testing"
)

type LearningRepositoryStub struct {
  called bool
  learning Learning
}

func (self *LearningRepositoryStub) StoreLearning(ctx context.Context, learning Learning)(string, error) {
  self.called = true
  self.learning = learning
  return self.learning.Id, nil
}

func TestLearningService_FailsEarly(t *testing.T) {
  repositoryStub := LearningRepositoryStub{false, Learning{}}
  ls := NewService(&repositoryStub)
  invalidLearning := Learning{}
  id, err := ls.StoreLearning(context.Background(), invalidLearning)
  if len(id) > 0 || err == nil {
    t.Errorf("Expected learning to be invalid: %v, but got id: %s, err: %v", invalidLearning, id, err)
  }

  if repositoryStub.called == true {
    t.Errorf("Expected learning repository to not be called, got %v", repositoryStub)
  }
}

func TestLearningService_ReturnsUnderlyingId(t *testing.T) {
  repositoryStub := LearningRepositoryStub{false, Learning{}}
  ls := LearningService{&repositoryStub}
  learningId := "someLearningId"

  learning := Learning{
    Id: learningId,
    Description: "SomeDescription",
    Topics: []string{"lol", "kek"},
    ReporterId: "SomeReporterId",
    Confirmed: false,
  }

  id, err := ls.StoreLearning(context.Background(), learning)
  if err != nil {
    t.Errorf("Expected storing valid learning to yield learning id, got %v", err)
  }
  if id != learningId {
    t.Errorf("Expected new learning id to be 'someLearningId' got %s", id)
  }
}
