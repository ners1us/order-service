package service

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) (*models.User, error)
	Login(email, password string) (string, error)
	DummyLogin(role string) (string, error)
}

type userServiceImpl struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
}

func NewUserService(userRepo repositories.UserRepository, jwtService JWTService) UserService {
	return &userServiceImpl{
		userRepo,
		jwtService,
	}
}

func (us *userServiceImpl) Register(user *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &models.User{}, err
	}
	user.Password = string(hashedPassword)
	user.ID = uuid.New().String()
	err = us.userRepo.CreateUser(user)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (us *userServiceImpl) Login(email, password string) (string, error) {
	user, err := us.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", enums.ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", enums.ErrWrongCredentials
	}
	return us.jwtService.GenerateToken(user.ID, user.Role)
}

func (us *userServiceImpl) DummyLogin(role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", enums.ErrInvalidRole
	}
	dummyUserID := "dummy_" + role
	return us.jwtService.GenerateToken(dummyUserID, role)
}
