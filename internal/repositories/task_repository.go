package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/models"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrDuplicateTitle = errors.New("duplicate task title")
	ErrInvalidFilter  = errors.New("invalid filter parameter")
)

type TaskRepository interface {
	GetTasksByUserID(userID int, categoryID string, status string, priority string) ([]models.Task, error)
	GetTaskByID(id int, userID int) (*models.Task, error)
	GetTaskByTitleAndUserID(title string, userID int) (*models.Task, error)
	CreateTask(task *models.Task) (*models.Task, error)
	UpdateTask(task *models.Task) (*models.Task, error)
	DeleteTask(id int, userID int) error
	ToggleTaskStatus(id int, userID int) (*models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) GetTasksByUserID(userID int, categoryID string, status string, priority string) ([]models.Task, error) {
	query := tr.db.Model(&models.Task{}).
		Where("user_id = ?", userID).
		Preload("Category").
		Preload("User")

	if categoryID != "" {
		if _, err := strconv.Atoi(categoryID); err != nil {
			return nil, ErrInvalidFilter
		}
		query = query.Where("category_id = ?", categoryID)
	}

	if status != "" {
		if _, err := strconv.ParseBool(status); err != nil {
			return nil, ErrInvalidFilter
		}
		query = query.Where("status = ?", status)
	}

	if priority != "" {
		priority = strings.ToLower(priority)
		if !models.IsValidPriority(models.Priority(priority)) {
			return nil, ErrInvalidFilter
		}
		query = query.Where("priority = ?", priority)
	}

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("error fetching tasks: %w", err)
	}

	return tasks, nil
}

func (tr *taskRepository) GetTaskByID(id int, userID int) (*models.Task, error) {
	var task models.Task
	err := tr.db.
		Preload("Category").
		Preload("User").
		Where("id = ? AND user_id = ?", id, userID).
		First(&task).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}

	return &task, err
}

func (tr *taskRepository) GetTaskByTitleAndUserID(title string, userID int) (*models.Task, error) {
	var task models.Task
	err := tr.db.
		Where("user_id = ? AND title = ?", userID, title).
		First(&task).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &task, err
}

func (tr *taskRepository) CreateTask(task *models.Task) (*models.Task, error) {
	err := tr.db.Create(task).Error
	if err != nil {
		if utils.IsDuplicateError(err) {
			return nil, ErrDuplicateTitle
		}
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	return tr.GetTaskByID(int(task.ID), int(task.UserID))
}

func (tr *taskRepository) UpdateTask(task *models.Task) (*models.Task, error) {
	err := tr.db.Save(task).Error
	if err != nil {
		if utils.IsDuplicateError(err) {
			return nil, ErrDuplicateTitle
		}
		return nil, fmt.Errorf("error updating task: %w", err)
	}

	return tr.GetTaskByID(int(task.ID), int(task.UserID))
}

func (tr *taskRepository) DeleteTask(id int, userID int) error {
	result := tr.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Task{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (tr *taskRepository) ToggleTaskStatus(id int, userID int) (*models.Task, error) {
	task, err := tr.GetTaskByID(id, userID)
	if err != nil {
		return nil, err
	}

	task.Status = !task.Status
	return tr.UpdateTask(task)
}
