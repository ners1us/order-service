package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	userID := "222"
	role := "employee"

	// Act
	token, err := jwtService.GenerateToken(userID, role)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	userID := "444"
	role := "employee"
	token, _ := jwtService.GenerateToken(userID, role)

	// Act
	claims, err := jwtService.ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	secretKey := "secret_for_testing"
	jwtService := NewJWTService(secretKey)
	invalidToken := "invalid_token_str"

	// Act
	claims, err := jwtService.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}
