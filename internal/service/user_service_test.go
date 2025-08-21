package service

import (
	"card_manage/internal/model"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository is a mock implementation of the IUserRepository interface.
type mockUserRepository struct {
	CreateUserFunc     func(user *model.User) (int64, error)
	GetUserByEmailFunc func(email string) (*model.User, error)
	GetUserByIDFunc    func(id int64) (*model.User, error)
	UpdateUserFunc     func(user *model.User) error
	DeleteUserFunc     func(id int64) error
}

// CreateUser delegates the call to the mock function.
func (m *mockUserRepository) CreateUser(user *model.User) (int64, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return 0, errors.New("CreateUserFunc not implemented")
}

// GetUserByEmail delegates the call to the mock function.
func (m *mockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, errors.New("GetUserByEmailFunc not implemented")
}

// GetUserByID delegates the call to the mock function.
func (m *mockUserRepository) GetUserByID(id int64) (*model.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, errors.New("GetUserByIDFunc not implemented")
}

// UpdateUser delegates the call to the mock function.
func (m *mockUserRepository) UpdateUser(user *model.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(user)
	}
	return errors.New("UpdateUserFunc not implemented")
}

// DeleteUser delegates the call to the mock function.
func (m *mockUserRepository) DeleteUser(id int64) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return errors.New("DeleteUserFunc not implemented")
}

func TestUserService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows // Simulate user not existing
			},
			CreateUserFunc: func(user *model.User) (int64, error) {
				return 1, nil // Simulate successful creation
			},
		}
		userService := NewUserService(mockRepo)

		user, err := userService.Register("test@example.com", "password123", "PLAYER")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		if user != nil {
			assert.Equal(t, int64(1), user.ID)
			assert.Equal(t, "test@example.com", user.Email)
		}
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*model.User, error) {
				return &model.User{}, nil // Simulate user already exists
			},
		}
		userService := NewUserService(mockRepo)

		_, err := userService.Register("test@example.com", "password123", "PLAYER")

		assert.Error(t, err)
		assert.Equal(t, ErrEmailExists, err)
	})
}

func TestUserService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         "PLAYER",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*model.User, error) {
				if email == mockUser.Email {
					return mockUser, nil
				}
				return nil, sql.ErrNoRows
			},
		}
		userService := NewUserService(mockRepo)

		user, err := userService.Login("test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.ID, user.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
		}
		userService := NewUserService(mockRepo)

		_, err := userService.Login("nonexistent@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByEmailFunc: func(email string) (*model.User, error) {
				return mockUser, nil
			},
		}
		userService := NewUserService(mockRepo)

		_, err := userService.Login("test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPassword, err)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	mockUser := &model.User{ID: 1, Email: "test@example.com"}

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByIDFunc: func(id int64) (*model.User, error) {
				if id == mockUser.ID {
					return mockUser, nil
				}
				return nil, sql.ErrNoRows
			},
		}
		userService := NewUserService(mockRepo)

		user, err := userService.GetUserByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.ID, user.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			GetUserByIDFunc: func(id int64) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
		}
		userService := NewUserService(mockRepo)

		_, err := userService.GetUserByID(2)

		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})
}