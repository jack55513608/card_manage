package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrDatabase       = errors.New("database error")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// UserService provides user-related services.
type UserService struct {
	userRepo repository.IUserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register creates a new user.
func (s *UserService) Register(email, password, role string) (*model.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	id, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	user.ID = id
	return user, nil
}

// Login validates user credentials and returns the user if successful.
func (s *UserService) Login(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	return user, nil
}
