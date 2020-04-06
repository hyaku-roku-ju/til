package learning
import (
  "testing"
)

func TestLearningParallel(t *testing.T) {
  description := "Some Description"
  topics := []string{"memes"}
  reporterId := "someReporterId"
  confirmed := false

  t.Run("IsValid", func(t *testing.T) {
    t.Run("Needs description", func(t *testing.T) {
      learning := Learning{"", topics, reporterId, confirmed}
      valid, err := learning.IsValid()
      if valid == true || err == nil {
        t.Errorf("Expected %v to be invalid, got valid: %v, %s", learning, valid, err)
      }
    })
    t.Run("Needs topics", func(t *testing.T) {
      learning := Learning{description, []string{}, reporterId, confirmed}
      valid, err := learning.IsValid()
      if valid == true || err == nil {
        t.Errorf("Expected %v to be invalid, got valid: %v, %s", learning, valid, err)
      }
    })
    t.Run("Topics must be unique", func(t *testing.T) {
      learning := Learning{description, []string{"memes", "memes"}, reporterId, confirmed}
      valid, err := learning.IsValid()
      if valid == true || err == nil {
        t.Errorf("Expected %v to be invalid, got valid: %v, %s", learning, valid, err)
      }
    })
    t.Run("Needs reporterId", func(t *testing.T) {
      learning := Learning{description, topics, "", confirmed}
      valid, err := learning.IsValid()
      if valid == true || err == nil {
        t.Errorf("Expected %v to be invalid, got valid: %v, %s", learning, valid, err)
      }
    })
  })
}
