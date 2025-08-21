package repository

import (
	"card_manage/internal/model"
	"database/sql"
)

// IUserRepository defines the interface for user repository operations.
type IUserRepository interface {
	CreateUser(user *model.User) (int64, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

// Statically check that UserRepository implements IUserRepository.
var _ IUserRepository = (*UserRepository)(nil)

// UserRepository handles database operations for users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database.
func (r *UserRepository) CreateUser(user *model.User) (int64, error) {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user's information in the database.
func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, role = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Role, user.ID)
	return err
}

// DeleteUser removes a user from the database.
func (r *UserRepository) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
