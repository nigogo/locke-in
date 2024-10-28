package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nigogo/locke-in/components"
	"github.com/nigogo/locke-in/renderer"
	"github.com/nigogo/locke-in/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		res := renderer.New(c.Request.Context(), http.StatusOK, views.GoalForm())
		c.Render(http.StatusOK, res)
	})

	r.POST("/goal", func(c *gin.Context) {
		println("Creating a new goal")

		var goal services.Goal
		if err := c.ShouldBind(&goal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		goalID := uuid.New().String()

		goalJSON, err := json.Marshal(goal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db[goalID] = string(goalJSON)

		c.Redirect(http.StatusSeeOther, "/goal/"+goalID)
	})

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

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goal))
		c.Render(http.StatusOK, res)
	})

	r.GET("/goals", func(c *gin.Context) {
		log.Println("Getting all goals")

		var goals []services.Goal
		for _, goalJSON := range db {
			var goal services.Goal
			if err := json.Unmarshal([]byte(goalJSON), &goal); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			goals = append(goals, goal)
		}

		//log all goals
		for _, goal := range goals {
			log.Println(goal)
		}

		res := renderer.New(c.Request.Context(), http.StatusOK, views.Goal(goals[0]))
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
