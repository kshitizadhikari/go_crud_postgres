package user

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"go_crud_postgres/internal/models"
	"go_crud_postgres/pkg/storage"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo    *UserRepository
	storage *storage.MinIO
}

func NewUserService(repo *UserRepository, storage *storage.MinIO) *UserService {
	return &UserService{repo: repo, storage: storage}
}

type CreateUserRequest struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest, fileHeader *multipart.FileHeader) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}

	defer file.Close()

	objectKey := fmt.Sprintf("avatars/%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	err = s.storage.Upload(
		ctx,
		objectKey,
		file,
		fileHeader.Size,
		fileHeader.Header.Get("Content-Type"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	// Create user
	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		AvatarKey: objectKey,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		_ = s.storage.Delete(ctx, objectKey)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (s *UserService) GetById(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*models.User, error) {
	user, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != user.Email {
		exists, err := s.repo.Exists(req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return nil, errors.New("email already taken by another user")
		}
		user.Email = req.Email
	}

	user.Name = req.Name

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err.Error())
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
