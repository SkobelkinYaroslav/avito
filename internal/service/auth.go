package service

import (
	"avito/internal/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	AuthRepository AuthRepository
}

func NewAuthService(authRepository AuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}

func (s *AuthService) RegisterService(req domain.AuthStruct) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPassword)

	if err := s.AuthRepository.RegisterRepo(req); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) LoginService(req domain.AuthStruct) (string, error) {
	user, err := s.AuthRepository.GetUserRepo(req)

	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token, err := s.createJWTService(user)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *AuthService) CheckTokenService(tokenString string) (domain.AuthStruct, error) {
	if tokenString == "" {
		return domain.AuthStruct{}, fmt.Errorf("empty token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return domain.AuthStruct{}, fmt.Errorf("error while parsing token")
	}

	var user domain.AuthStruct

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.AuthStruct{}, fmt.Errorf("token expired")
		}

		req := domain.AuthStruct{
			Email: claims["sub"].(string),
		}
		user, err = s.AuthRepository.GetUserRepo(req)

		if err != nil {
			return domain.AuthStruct{}, fmt.Errorf("error while getting user")
		}
	} else {
		return domain.AuthStruct{}, fmt.Errorf("error while getting claims")
	}

	return user, nil
}

func (s *AuthService) createJWTService(user domain.AuthStruct) (string, error) {
	user, err := s.AuthRepository.GetUserRepo(user)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
