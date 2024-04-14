package service_test

import (
	"avito/internal/domain"

	mocks "avito/internal/mocks/repository"
	"avito/internal/service"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserBannerService(t *testing.T) {
	mockBannerRepo := &mocks.BannerRepository{}

	s := service.NewBannerService(mockBannerRepo)

	req := domain.GetUserBannerRequest{}
	expectedBanner := domain.Banner{}

	mockBannerRepo.On("GetUserBannerRepo", req).Return(expectedBanner, nil)

	banner, err := s.GetUserBannerService(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedBanner, banner)
	mockBannerRepo.AssertExpectations(t)
}

func TestGetAllBannersService(t *testing.T) {
	mockBannerRepo := &mocks.BannerRepository{}

	s := service.NewBannerService(mockBannerRepo)

	req := domain.GetBannersRequest{}
	expectedBanners := []domain.Banner{{}}

	mockBannerRepo.On("GetAllBannersRepo", req).Return(expectedBanners, nil)

	banners, err := s.GetAllBannersService(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedBanners, banners)
	mockBannerRepo.AssertExpectations(t)
}

func TestPostBannerService(t *testing.T) {
	mockBannerRepo := &mocks.BannerRepository{}

	s := service.NewBannerService(mockBannerRepo)

	req := domain.Banner{}
	expectedBanner := domain.Banner{}

	mockBannerRepo.On("GetUserBannerRepo", mock.Anything).Return(domain.Banner{}, sql.ErrNoRows)
	mockBannerRepo.On("PostBannerRepo", req).Return(expectedBanner, nil)

	banner, err := s.PostBannerService(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedBanner, banner)
	mockBannerRepo.AssertExpectations(t)
}

func TestPatchBannerService(t *testing.T) {
	mockBannerRepo := &mocks.BannerRepository{}

	s := service.NewBannerService(mockBannerRepo)

	req := domain.Banner{}

	mockBannerRepo.On("PatchBannerRepo", req).Return(nil)

	err := s.PatchBannerService(req)

	assert.NoError(t, err)
	mockBannerRepo.AssertExpectations(t)
}

func TestDeleteBannerService(t *testing.T) {
	mockBannerRepo := &mocks.BannerRepository{}

	s := service.NewBannerService(mockBannerRepo)

	req := domain.Banner{
		BannerID: 1,
	}

	mockBannerRepo.On("DeleteBannerRepo", req).Return(nil)

	err := s.DeleteBannerService(req)

	assert.NoError(t, err)
	mockBannerRepo.AssertExpectations(t)
}
