package service_test

import (
	"avito/internal/domain"
	mocks "avito/internal/mocks/repository"
	"avito/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
	"time"
)

func TestRegisterService(t *testing.T) {
	mockAuthRepo := &mocks.AuthRepository{}

	s := service.NewAuthService(mockAuthRepo)

	expectedUser := domain.AuthStruct{
		Email:    "test@test.com",
		Password: "password",
		IsAdmin:  false,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(expectedUser.Password), bcrypt.DefaultCost)
	expectedUser.Password = string(hashedPassword)

	mockAuthRepo.On("RegisterRepo", mock.Anything).Return(nil)

	err := s.RegisterService(expectedUser)

	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestLoginService(t *testing.T) {
	mockAuthRepo := &mocks.AuthRepository{}

	s := service.NewAuthService(mockAuthRepo)

	expectedUser := domain.AuthStruct{
		Email:    "test@test.com",
		Password: "password",
		IsAdmin:  false,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(expectedUser.Password), bcrypt.DefaultCost)
	hashedUser := expectedUser
	hashedUser.Password = string(hashedPassword)

	mockAuthRepo.On("GetUserRepo", mock.Anything).Return(hashedUser, nil)

	token, err := s.LoginService(expectedUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockAuthRepo.AssertExpectations(t)
}

func TestCheckTokenService(t *testing.T) {
	mockAuthRepo := &mocks.AuthRepository{}

	s := service.NewAuthService(mockAuthRepo)

	expectedUser := domain.AuthStruct{
		Email:    "test@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	hashedUser := expectedUser
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(expectedUser.Password), bcrypt.DefaultCost)
	hashedUser.Password = string(hashedPassword)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": expectedUser.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))

	expectedUser.Password = ""
	mockAuthRepo.On("GetUserRepo", expectedUser).Return(hashedUser, nil)

	user, err := s.CheckTokenService(tokenString)

	assert.NoError(t, err)
	assert.Equal(t, hashedUser, user)
	mockAuthRepo.AssertExpectations(t)
}
