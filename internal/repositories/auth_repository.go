package repositories

import (
	"errors"

	"github.com/alastor-4/sylcot-go-gin-backend/internal/models"
	"github.com/alastor-4/sylcot-go-gin-backend/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrTokenNotFound     = errors.New("token not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type AuthRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByToken(token string) (*models.User, error)
	FindByResetToken(token string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		if utils.IsDuplicateError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *authRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
