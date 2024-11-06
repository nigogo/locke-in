package services

import (
	"fmt"
	"time"

	"github.com/nigogo/locke-in/db"
)

type Goal struct {
	ID        string    `form:"id"`
	Name      string    `form:"name" binding:"required"`
	StartDate time.Time `form:"endDate" time_format:"2006-01-02T15:04"`
	EndDate   time.Time `form:"endDate" binding:"required" time_format:"2006-01-02T15:04"`
	Completed bool      `form:"completed"`
}

func GetGoals() ([]Goal, error) {
	var goals []Goal
	result := db.GetDB().Find(&goals)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get goals: %w", result.Error)
	}
	return goals, nil
}

func GetGoal(id string) (Goal, error) {
	var goal Goal
	result := db.GetDB().First(&goal, "id = ?", id)
	if result.Error != nil {
		return Goal{}, fmt.Errorf("failed to get goal: %w", result.Error)
	}
	return goal, nil
}

func GetActiveGoal() (*Goal, error) {
	var goal Goal
	result := db.GetDB().First(&goal, "completed = ?", false)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active goal: %w", result.Error)
	}
	return &goal, nil
}

func GetCompletedGoals() ([]Goal, error) {
	var goals []Goal
	result := db.GetDB().Find(&goals, "completed = ?", true)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get completed goals: %w", result.Error)
	}
	return goals, nil
}

func StoreGoal(goal Goal) error {
	result := db.GetDB().Save(&goal)
	if result.Error != nil {
		return fmt.Errorf("failed to store goal: %w", result.Error)
	}
	return nil
}
