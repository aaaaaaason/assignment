package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	code := http.StatusInternalServerError
	errMsg := "unknown error"

	if custom, ok := c.Errors[0].Err.(*Error); ok {
		code = custom.Code()
		errMsg = custom.Error()
	}

	c.JSON(code, gin.H{"error": errMsg})
}

func setupRouter(controller *TaskController) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(ErrorHandlingMiddleware)

	r.GET("/tasks", controller.List)
	r.POST("/tasks", controller.Create)
	r.PUT("/tasks/:id", controller.Update)
	r.DELETE("/tasks/:id", controller.Delete)
	return r
}

func main() {
	controller := &TaskController{
		database: make(map[string]*Task),
	}

	r := setupRouter(controller)
	r.Run("0.0.0.0:8080")
}
