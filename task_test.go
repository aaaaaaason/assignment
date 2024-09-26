package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTaskController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TaskController Suite")
}

var _ = Describe("Tasks", func() {
	var (
		router     *gin.Engine
		controller TaskController
		recorder   *httptest.ResponseRecorder
		req        *http.Request
		respBody   []byte
		err        error
		tasks      []Task
	)

	BeforeEach(func() {
		controller.database = make(map[string]*Task)
		router = setupRouter(&controller)
		recorder = httptest.NewRecorder()

		tasks = []Task{
			{
				ID:     uuid.New().String(),
				Name:   "task1",
				Status: Incompleted,
			},
			{
				ID:     uuid.New().String(),
				Name:   "task2",
				Status: Completed,
			},
			{
				ID:     uuid.New().String(),
				Name:   "task3",
				Status: Completed,
			},
		}

		for i := range tasks {
			controller.database[tasks[i].ID] = &tasks[i]
		}
	})

	JustBeforeEach(func() {
		router.ServeHTTP(recorder, req)
		respBody = recorder.Body.Bytes()
	})

	Describe("List Tasks", func() {
		BeforeEach(func() {
			req, err = http.NewRequest(http.MethodGet, "/tasks", nil)
			Expect(err).To(Succeed())
		})

		It("Succeed", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
			got := []Task{}
			err = json.Unmarshal(respBody, &got)
			Expect(err).To(Succeed())
			Expect(got).To(ConsistOf(tasks[0], tasks[1], tasks[2]))
		})
	})

	Describe("Create Task", func() {
		Context("Succeed", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "task4",
					Status: Completed,
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Succeed", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))

				var got Task
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())

				createdTask := controller.database[got.ID]
				Expect(createdTask.Name).To(Equal(task.Name))
				Expect(createdTask.Status).To(Equal(task.Status))
			})
		})

		Context("Empty task name", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "",
					Status: Completed,
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				got := map[string]string{}
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())
				Expect(got["error"]).To(Equal("empty task name"))
			})
		})

		Context("Invalid task status", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "task4",
					Status: TaskStatus(5),
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				got := map[string]string{}
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())
				Expect(got["error"]).To(Equal("invalid task status"))
			})
		})
	})

	Describe("Update Task", func() {
		Context("Succeed", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "task1 - updated",
					Status: Completed,
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPut, "/tasks/"+tasks[0].ID, bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Succeed", func() {
				Expect(recorder.Code).To(Equal(http.StatusNoContent))

				updatedTask := controller.database[tasks[0].ID]
				Expect(updatedTask.Name).To(Equal(task.Name))
				Expect(updatedTask.Status).To(Equal(task.Status))
			})
		})

		Context("Task not found", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "task1 - updated",
					Status: Completed,
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPut, "/tasks/"+uuid.New().String(), bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusNotFound))

				got := map[string]string{}
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())
				Expect(got["error"]).To(Equal("task not found"))
			})
		})

		Context("Empty task name", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "",
					Status: Completed,
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPut, "/tasks/"+tasks[0].ID, bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				got := map[string]string{}
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())
				Expect(got["error"]).To(Equal("empty task name"))
			})
		})

		Context("Invalid task status", func() {
			var task Task

			BeforeEach(func() {
				task = Task{
					Name:   "task1 - updated",
					Status: TaskStatus(5),
				}

				payload, err := json.Marshal(task)
				Expect(err).To(Succeed())

				req, err = http.NewRequest(http.MethodPut, "/tasks/"+tasks[0].ID, bytes.NewReader(payload))
				Expect(err).To(Succeed())
			})

			It("Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				got := map[string]string{}
				err = json.Unmarshal(respBody, &got)
				Expect(err).To(Succeed())
				Expect(got["error"]).To(Equal("invalid task status"))
			})
		})
	})

	Describe("Delete Task", func() {
		Context("Succeed", func() {
			BeforeEach(func() {
				req, err = http.NewRequest(http.MethodDelete, "/tasks/"+tasks[0].ID, nil)
				Expect(err).To(Succeed())
			})

			It("Succeed", func() {
				Expect(recorder.Code).To(Equal(http.StatusNoContent))

				_, found := controller.database[tasks[0].ID]
				Expect(found).To(BeFalse())
			})
		})

		Context("Task not found", func() {
			BeforeEach(func() {
				req, err = http.NewRequest(http.MethodDelete, "/tasks/"+uuid.New().String(), nil)
				Expect(err).To(Succeed())
			})

			It("Nothing is removed", func() {
				Expect(recorder.Code).To(Equal(http.StatusNoContent))
				Expect(len(controller.database)).To(Equal(3))
			})
		})
	})
})
