package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alastor-4/sylcot-go-gin-backend/internal/models"
	"github.com/alastor-4/sylcot-go-gin-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	repo repositories.TaskRepository
}

func NewTaskHandler(repo repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// GetTasks godoc
// @Summary Get filtered tasks
// @Description Get tasks with optional filters for category, status, and priority
// @Tags tasks
// @Produce json
// @Param categoryId query int false "Filter by category ID"
// @Param status query boolean false "Filter by completion status (true/false)"
// @Param priority query string false "Filter by priority (high/medium/low)" Enums(high, medium, low)
// @Security ApiKeyAuth
// @Success 200 {array} models.Task
// @Failure 500 {object} object{error=string}
// @Router /api/tasks [get]
func (th *TaskHandler) GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	categoryID := c.Query("categoryId")
	status := c.Query("status")
	priority := c.Query("priority")

	tasks, err := th.repo.GetTasksByUserID(userID.(int), categoryID, status, priority)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidFilter) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.TaskRequest true "Task creation data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Task
// @Failure 400 {object} object{error=string,details=object}
// @Failure 500 {object} object{error=string}
// @Router /api/tasks [post]
func (th *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq models.TaskRequest
	userID, _ := c.Get("userID")

	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	if err := models.ValidateTaskRequest(taskReq); err != nil {
		validationErrors := models.GetTaskValidationMessages(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	existingTask, err := th.repo.GetTaskByTitleAndUserID(taskReq.Title, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking task existence"})
		return
	}
	if existingTask != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists."})
		return
	}

	task := models.Task{
		Title:      taskReq.Title,
		Priority:   taskReq.Priority,
		CategoryID: taskReq.CategoryID,
		UserID:     uint(userID.(int)),
	}

	newTask, err := th.repo.CreateTask(&task)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task title already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task's details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.TaskRequest true "Task update data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Task
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/tasks/{id} [put]
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	existingTask, err := th.repo.GetTaskByID(id, userID.(int))
	if err != nil {
		if errors.Is(err, repositories.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching task"})
		return
	}

	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	if err := models.ValidateTaskRequest(taskReq); err != nil {
		validationErrors := models.GetTaskValidationMessages(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	if existingTask.Title == taskReq.Title &&
		existingTask.Priority == taskReq.Priority &&
		existingTask.CategoryID == taskReq.CategoryID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes detected"})
		return
	}

	existingTask.Title = taskReq.Title
	existingTask.Priority = taskReq.Priority
	existingTask.CategoryID = taskReq.CategoryID

	updatedTask, err := th.repo.UpdateTask(existingTask)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task title already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Permanently delete a task
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/tasks/{id} [delete]
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	err := th.repo.DeleteTask(id, userID.(int))
	if err != nil {
		if errors.Is(err, repositories.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Task %d deleted", id)})
}

// ToggleTask godoc
// @Summary Toggle task status
// @Description Toggle a task's completion status
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Task
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/tasks/{id}/complete [patch]
func (th *TaskHandler) ToggleTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	task, err := th.repo.ToggleTaskStatus(id, userID.(int))
	if err != nil {
		if errors.Is(err, repositories.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error toggling task status"})
		return
	}

	c.JSON(http.StatusOK, task)
}
