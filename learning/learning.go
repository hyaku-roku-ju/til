package learning

import (
  "fmt"
)

type Learning struct {
  Description string
  Topics []string
  ReporterId string
  Confirmed bool
}

func removeDuplicateStrings(list []string) []string {
  set := make(map[string]bool)
  withoutDuplicates := make([]string, 0)
  
  for _, element := range list {
    _, exists := set[element]
    if exists {
      continue
    }
    set[element] = true
    withoutDuplicates = append(withoutDuplicates,element)
  }
  
  return withoutDuplicates
}

func (learning *Learning) IsValid() (bool, error) {
  if len(learning.Description) == 0 {
    return false, fmt.Errorf("Description of learning cannot be empty")
  }
  if len(learning.Topics) == 0 {
    return false, fmt.Errorf("Learning must have at least one topic, got %v", learning.Topics)
  }
  if len(learning.Topics) != len(removeDuplicateStrings(learning.Topics)) {
    return false, fmt.Errorf("Learning contains duplicate topics, %v", learning.Topics)
  }
  if len(learning.ReporterId) == 0 {
    return false, fmt.Errorf("Learning must have a reporterId associated with it")
  }
  return true, nil
} 

func New(reporterId string, topics []string, description string) (Learning, error) {
  learning := Learning{
    Description: description,
    Topics: topics,
    ReporterId: reporterId,
    Confirmed: false,
  }    

  if _, err := learning.IsValid(); err != nil {
    return Learning{}, err
  } else {
    return learning, nil
  }
}
