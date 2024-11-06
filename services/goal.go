package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Goal struct {
	ID        string    `form:"id"`
	Name      string    `form:"name" binding:"required"`
	StartDate time.Time `form:"endDate" time_format:"2006-01-02T15:04"`
	EndDate   time.Time `form:"endDate" binding:"required" time_format:"2006-01-02T15:04"`
	Completed bool      `form:"completed"`
}

var db = make(map[string]string)

// GetGoals parses all goals from db and returns a slice of goals and an error if parsing fails.
func GetGoals() ([]Goal, error) {
	var goals []Goal
	for _, goalJSON := range db {
		var goal Goal
		if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
			return nil, fmt.Errorf("failed to unmarshal goal: %w", err)
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func GetGoal(id string) (Goal, error) {
	goalJSON, ok := db[id]
	if !ok {
		return Goal{}, fmt.Errorf("goal not found")
	}
	var goal Goal
	if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
		return Goal{}, fmt.Errorf("failed to unmarshal goal: %w", err)
	}
	return goal, nil
}

// GetActiveGoal returns a pointer to the first active goal, or nil if none is found.
func GetActiveGoal() (*Goal, error) {
	goals, err := GetGoals()
	if err != nil {
		return nil, err
	}

	for i := range goals {
		if !goals[i].Completed {
			return &goals[i], nil
		}
	}
	return nil, nil // No active goal found
}

// GetCompletedGoals filters and returns all completed goals.
func GetCompletedGoals() ([]Goal, error) {
	goals, err := GetGoals()
	if err != nil {
		return nil, err
	}

	log.Printf("goals: %v", goals)

	completedGoals := make([]Goal, 0, len(goals))
	for _, goal := range goals {
		if goal.Completed {
			completedGoals = append(completedGoals, goal)
		}
	}

	log.Printf("completed goals: %v", completedGoals)

	// log all GetCompletedGoals
	for i, goal := range completedGoals {
		log.Printf("completed goal %d: %v", i, goal)
	}

	return completedGoals, nil
}

func StoreGoal(goal Goal) error {
	goalJson, err := json.Marshal(goal)
	if err != nil {
		return fmt.Errorf("failed to marshal goal: %w", err)
	}

	db[goal.ID] = string(goalJson)
	return nil
}
