package service

import (
	"avito/internal/domain"
)

type AuthRepository interface {
	RegisterRepo(user domain.AuthStruct) error
	GetUserRepo(user domain.AuthStruct) (domain.AuthStruct, error)
}

type BannerRepository interface {
	GetUserBannerRepo(req domain.GetUserBannerRequest) (domain.Banner, error)
	GetAllBannersRepo(req domain.GetBannersRequest) ([]domain.Banner, error)
	PostBannerRepo(req domain.Banner) (domain.Banner, error)
	PatchBannerRepo(req domain.Banner) error
	DeleteBannerRepo(req domain.Banner) error
}
