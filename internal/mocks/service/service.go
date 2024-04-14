package mocks

import (
	"avito/internal/domain"
	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func (m *AuthService) RegisterService(req domain.AuthStruct) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *AuthService) LoginService(req domain.AuthStruct) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *AuthService) CheckTokenService(token string) (domain.AuthStruct, error) {
	args := m.Called(token)
	return args.Get(0).(domain.AuthStruct), args.Error(1)
}

type BannerService struct {
	mock.Mock
}

func (m *BannerService) GetUserBannerService(req domain.GetUserBannerRequest) (domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).(domain.Banner), args.Error(1)
}

func (m *BannerService) GetAllBannersService(req domain.GetBannersRequest) ([]domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).([]domain.Banner), args.Error(1)
}

func (m *BannerService) PostBannerService(req domain.Banner) (domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).(domain.Banner), args.Error(1)
}

func (m *BannerService) PatchBannerService(req domain.Banner) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *BannerService) DeleteBannerService(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
