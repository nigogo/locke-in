package services

import "time"

type Goal struct {
	ID        string    `form:"id"`
	Name      string    `form:"name" binding:"required"`
	Deadline  time.Time `form:"deadline" binding:"required" time_format:"2006-01-02T15:04"`
	Completed bool      `form:"completed"`
}
