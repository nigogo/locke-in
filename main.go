package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	views "github.com/nigogo/locke-in/components"
	"github.com/nigogo/locke-in/renderer"
	"github.com/nigogo/locke-in/services"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

// GetGoals parses all goals from db and returns a slice of goals and an error if parsing fails.
func GetGoals() ([]services.Goal, error) {
	var goals []services.Goal
	for _, goalJSON := range db {
		var goal services.Goal
		if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
			return nil, fmt.Errorf("failed to unmarshal goal: %w", err)
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

// GetActiveGoal returns a pointer to the first active goal, or nil if none is found.
func GetActiveGoal() (*services.Goal, error) {
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
func GetCompletedGoals() ([]services.Goal, error) {
	goals, err := GetGoals()
	if err != nil {
		return nil, err
	}

	log.Printf("goals: %v", goals)

	completedGoals := make([]services.Goal, 0, len(goals))
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

func storeGoal(goal services.Goal) error {
	goalJson, err := json.Marshal(goal)
	if err != nil {
		return fmt.Errorf("failed to marshal goal: %w", err)
	}

	db[goal.ID] = string(goalJson)
	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		res := renderer.New(c.Request.Context(), http.StatusOK, views.GoalForm())
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

			goalJson, err := json.Marshal(goal)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			db[goalID] = string(goalJson)

			c.Redirect(http.StatusSeeOther, "/goal/"+goalID)
		},
	)

	r.GET("/goal/:id", func(c *gin.Context) {
		goalID := c.Param("id")
		goalJSON, ok := db[goalID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
			return
		}

		var goal services.Goal
		if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		completedGoals, err := GetCompletedGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goal, completedGoals))
		c.Render(http.StatusOK, res)
	})

	r.PATCH("/goal/:id", func(c *gin.Context) {
		goalID := c.Param("id")
		goalJSON, ok := db[goalID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
			return
		}

		var goal services.Goal
		if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		goal.Completed = true
		err := storeGoal(goal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store goal"})
			return
		}

		completedGoals, err := GetCompletedGoals()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goal, completedGoals))
		c.Render(http.StatusOK, res)
	})

	r.GET("/goals", func(c *gin.Context) {
		log.Println("Getting all goals")

		activeGoal, err := GetActiveGoal()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"could not get goals": err.Error()})
		}

		allGoals, err := GetGoals()
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

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar",
		"nico": "123",
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
	r.Static("/assets", "./assets")
	_ = r.Run(":8080")
}
