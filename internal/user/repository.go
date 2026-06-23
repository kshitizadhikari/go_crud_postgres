package user

import (
	"context"
	"go_crud_postgres/internal/database"

	"go_crud_postgres/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).Order("id").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
