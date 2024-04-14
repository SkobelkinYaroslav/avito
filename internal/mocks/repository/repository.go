package mocks

import (
	"avito/internal/domain"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	AuthRepository
	BannerRepository
}

type AuthRepository struct {
	mock.Mock
}

func (m *AuthRepository) RegisterRepo(user domain.AuthStruct) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *AuthRepository) GetUserRepo(user domain.AuthStruct) (domain.AuthStruct, error) {
	args := m.Called(user)
	return args.Get(0).(domain.AuthStruct), args.Error(1)
}

type BannerRepository struct {
	mock.Mock
}

func (m *BannerRepository) GetUserBannerRepo(req domain.GetUserBannerRequest) (domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).(domain.Banner), args.Error(1)
}

func (m *BannerRepository) GetAllBannersRepo(req domain.GetBannersRequest) ([]domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).([]domain.Banner), args.Error(1)
}

func (m *BannerRepository) PostBannerRepo(req domain.Banner) (domain.Banner, error) {
	args := m.Called(req)
	return args.Get(0).(domain.Banner), args.Error(1)
}

func (m *BannerRepository) PatchBannerRepo(req domain.Banner) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *BannerRepository) DeleteBannerRepo(req domain.Banner) error {
	args := m.Called(req)
	return args.Error(0)
}
