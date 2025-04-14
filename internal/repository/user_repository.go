package repository

import (
	"database/sql"
	"errors"
	"github.com/ners1us/order-service/internal/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (ur *userRepositoryImpl) CreateUser(user *model.User) error {
	_, err := ur.db.Exec("INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.Password, user.Role)
	return err
}

func (ur *userRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.db.QueryRow("SELECT id, email, password, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return &model.User{}, nil
	}
	return &user, err
}
