package handlers

import (
	"context"
	"github.com/nigogo/locke-in/services"
)

type GoalService interface {
	GetAll(c context.Context) ([]services.Goal, error)
	GetByID(c context.Context, id string) (services.Goal, error)
}

type DefaultHandler struct {
	GoalService GoalService
}

func New(gs GoalService) *DefaultHandler {
	return &DefaultHandler{
		GoalService: gs,
	}
}
