package usecase

import (
	"context"
	"errors"
	"testing"

	"auth_service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)

	uc := NewAccessUsecase(mockRepo, "test-secret", mockEmail)

	ctx := context.Background()
	testEmail := "test@example.com"
	testPassword := "strongpass123"
	testName := "Adil"

	mockRepo.On("GetByEmail", ctx, testEmail).Return(nil, errors.New("not found"))

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(int64(1), nil)

	mockEmail.On("SendWelcomeEmail", testEmail, testName).Return(nil)

	userID, err := uc.Register(ctx, testEmail, testPassword, testName)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), userID)

	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	uc := NewAccessUsecase(mockRepo, "test-secret", mockEmail)

	ctx := context.Background()

	existingUser := &domain.User{Email: "test@example.com"}
	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

	_, err := uc.Register(ctx, "test@example.com", "password", "Adil")

	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())

	mockRepo.AssertNotCalled(t, "Create")
	mockEmail.AssertNotCalled(t, "SendWelcomeEmail")
}
