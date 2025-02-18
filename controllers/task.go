package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alastor-4/sylcot-go-gin-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
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
func (tc *TaskController) GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	query := tc.DB.Where("user_id = ?", userID)

	if categoryID := c.Query("categoryId"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status == "true")
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	var tasks []models.Task
	if err := query.Preload("Category").Preload("User").Find(&tasks).Error; err != nil {
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
func (tc *TaskController) CreateTask(c *gin.Context) {
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

	var existingTask models.Task
	if err := tc.DB.Where("user_id = ? AND title = ?", userID, taskReq.Title).First(&existingTask).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists."})
		return
	}

	task := models.Task{
		Title:      taskReq.Title,
		Priority:   taskReq.Priority,
		CategoryID: taskReq.CategoryID,
		UserID:     uint(userID.(int)),
	}

	if err := tc.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	if err := tc.DB.Preload("Category").Preload("User").First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load task associations"})
		return
	}

	c.JSON(http.StatusCreated, task)
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
func (tc *TaskController) UpdateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	var task models.Task
	if err := tc.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
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

	if task.Title == taskReq.Title && task.Priority == taskReq.Priority && task.CategoryID == taskReq.CategoryID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes detected. Please modify at least one field."})
		return
	}

	task.Title = taskReq.Title
	task.Priority = taskReq.Priority
	task.CategoryID = taskReq.CategoryID

	if err := tc.DB.Save(&task).Error; err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "duplicate") || strings.Contains(errStr, "unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	if err := tc.DB.Preload("Category").Preload("User").First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load task associations"})
		return
	}

	c.JSON(http.StatusOK, task)
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
func (tc *TaskController) DeleteTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	result := tc.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete task"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Task with id %d was successfully deleted.", id)})
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
func (tc *TaskController) ToggleTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	var task models.Task
	if err := tc.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Status = !task.Status
	if err := tc.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update task"})
		return
	}

	if err := tc.DB.Preload("Category").Preload("User").First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load task associations"})
		return
	}

	c.JSON(http.StatusOK, task)
}
