package service

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) (*model.User, error)
	Login(email, password string) (string, error)
	DummyLogin(role string) (string, error)
}

type userServiceImpl struct {
	userRepo   repository.UserRepository
	jwtService JWTService
}

func NewUserService(userRepo repository.UserRepository, jwtService JWTService) UserService {
	return &userServiceImpl{
		userRepo,
		jwtService,
	}
}

func (us *userServiceImpl) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &model.User{}, err
	}
	user.Password = string(hashedPassword)
	user.ID = uuid.New().String()
	err = us.userRepo.CreateUser(user)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func (us *userServiceImpl) Login(email, password string) (string, error) {
	user, err := us.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", enum.ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", enum.ErrWrongPassword
	}
	return us.jwtService.GenerateToken(user.ID, user.Role)
}

func (us *userServiceImpl) DummyLogin(role string) (string, error) {
	if !enum.IsValidRole(enum.Role(role)) {
		return "", enum.ErrInvalidRole
	}
	dummyUserID := "dummy_" + role
	return us.jwtService.GenerateToken(dummyUserID, role)
}
