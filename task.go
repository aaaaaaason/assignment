package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskStatus can be 0 (imcompleted) or 1 (completed).
type TaskStatus int

const (
	Incompleted TaskStatus = 0
	Completed   TaskStatus = 1
)

var (
	ErrParseRequestPayload = NewError(http.StatusBadRequest, "decode request payload failed")
	ErrTaskNotFound        = NewError(http.StatusNotFound, "task not found")
	ErrEmptyTaskName       = NewError(http.StatusBadRequest, "empty task name")
	ErrInvalidTaskStatus   = NewError(http.StatusBadRequest, "invalid task status")
)

// Task is the internal representation of task.
type Task struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

// TaskController is the controller of task APIs.
type TaskController struct {
	// in-memory data storage
	database map[string]*Task
}

func (t *TaskController) List(c *gin.Context) {
	tasks := make([]Task, 0, len(t.database))
	for _, v := range t.database {
		tasks = append(tasks, *v)
	}

	c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) Create(c *gin.Context) {
	task, err := t.validateTask(c)
	if err != nil {
		c.Error(err)
		return
	}

	task.ID = uuid.New().String()
	t.database[task.ID] = task

	c.JSON(http.StatusOK, gin.H{
		"id": task.ID,
	})
}

func (t *TaskController) Update(c *gin.Context) {
	id := c.Param("id")
	if _, ok := t.database[id]; !ok {
		c.Error(ErrTaskNotFound)
		return
	}

	task, err := t.validateTask(c)
	if err != nil {
		c.Error(err)
		return
	}

	task.ID = id
	t.database[task.ID] = task

	c.JSON(http.StatusNoContent, nil)
}

func (t *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, ok := t.database[id]; ok {
		delete(t.database, id)
	}

	c.JSON(http.StatusNoContent, nil)
}

func (t *TaskController) validateTask(c *gin.Context) (*Task, error) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		return nil, ErrParseRequestPayload
	}

	if len(task.Name) == 0 {
		return nil, ErrEmptyTaskName
	}

	if int(task.Status) != 0 && int(task.Status) != 1 {
		return nil, ErrInvalidTaskStatus
	}

	return &task, nil
}
