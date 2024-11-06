package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	views "github.com/nigogo/locke-in/components"
	services "github.com/nigogo/locke-in/services"

	"github.com/nigogo/locke-in/renderer"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		completedGoals, err := services.GetCompletedGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.GoalForm(completedGoals))
		c.Render(http.StatusOK, res)
	})

	r.POST(
		"/goal",
		func(c *gin.Context) {
			println("Creating a new goal")
			println("NOW: " + time.Now().String())

			var goal services.Goal
			if err := c.ShouldBind(&goal); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			println("local time: " + goal.EndDate.String())

			goalID := uuid.New().String()
			startDate := time.Now()

			goal = services.Goal{
				ID:        goalID,
				Name:      goal.Name,
				StartDate: startDate,
				EndDate:   goal.EndDate,
				Completed: false,
			}

			services.StoreGoal(goal)

			c.Redirect(http.StatusSeeOther, "/goal/"+goalID)
		},
	)

	r.GET("/goal/:id", func(c *gin.Context) {
		goal, err := services.GetGoal(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		}

		completedGoals, err := services.GetCompletedGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goal, completedGoals))
		c.Render(http.StatusOK, res)
	})

	r.PATCH("/goal/:id", func(c *gin.Context) {
		goal, err := services.GetGoal(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		}

		goal.Completed = true
		err = services.StoreGoal(goal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store goal"})
			return
		}

		completedGoals, err := services.GetCompletedGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goal, completedGoals))
		c.Render(http.StatusOK, res)
	})

	r.GET("/goals", func(c *gin.Context) {
		log.Println("Getting all goals")

		activeGoal, err := services.GetActiveGoal()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		var allGoals []services.Goal
		allGoals, err = services.GetGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(*activeGoal, allGoals))
		c.Render(http.StatusOK, res)
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong foo")
	})

	return r
}

func main() {
	db, err := gorm.Open(sqlite.Open("locke-in.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&services.Goal{})

	r := setupRouter()
	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
	r.Static("/assets", "./assets")
	_ = r.Run(":8080")
}
