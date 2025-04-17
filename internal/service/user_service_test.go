package service

import (
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestDummyLogin_InvalidRole(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	mockUserRepo := new(repository.MockUserRepository)
	service := NewUserService(mockUserRepo, jwtService)
	role := "programmer"

	// Act
	_, err := service.DummyLogin(role)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrInvalidRole, err)
}

func TestDummyLogin_Moderator(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	mockUserRepo := new(repository.MockUserRepository)
	service := NewUserService(mockUserRepo, jwtService)
	role := "moderator"

	// Act
	token, err := service.DummyLogin(role)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestRegister_Success(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	mockUserRepo := new(repository.MockUserRepository)
	service := NewUserService(mockUserRepo, jwtService)
	user := &model.User{Email: "tonyStark@example.com", Password: "password"}
	mockUserRepo.On("CreateUser", mock.Anything).Return(nil)

	// Act
	result, err := service.Register(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.NotEqual(t, "password", result.Password)
}

func TestLogin_WrongPassword(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	mockUserRepo := new(repository.MockUserRepository)
	service := NewUserService(mockUserRepo, jwtService)
	email := "coolmail@example.com"
	password := "password12345"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	user := &model.User{ID: "user1", Email: email, Password: string(hashedPassword), Role: "employee"}
	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)

	// Act
	_, err := service.Login(email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrWrongPassword, err)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	mockUserRepo := new(repository.MockUserRepository)
	service := NewUserService(mockUserRepo, jwtService)
	email := "linuxuser@mail.com"
	password := "password"
	mockUserRepo.On("GetUserByEmail", email).Return(&model.User{}, nil)

	// Act
	_, err := service.Login(email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrUserNotFound, err)
}
