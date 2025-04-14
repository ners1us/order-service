package repository

import (
	"database/sql"
	"errors"
	"github.com/ners1us/order-service/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (ur *userRepositoryImpl) CreateUser(user *models.User) error {
	_, err := ur.db.Exec("INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.Password, user.Role)
	return err
}

func (ur *userRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT id, email, password, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.User{}, nil
	}
	return &user, err
}
