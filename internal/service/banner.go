package service

import (
	"avito/internal/domain"
	errGroup "avito/internal/error"
	"database/sql"
	"log"
)

type BannerService struct {
	BannerRepository BannerRepository
}

func NewBannerService(bannerRepository BannerRepository) *BannerService {
	return &BannerService{
		BannerRepository: bannerRepository,
	}
}

func (s *BannerService) GetUserBannerService(req domain.GetUserBannerRequest) (domain.Banner, error) {
	banner, err := s.BannerRepository.GetUserBannerRepo(req)
	if err == sql.ErrNoRows {
		return domain.Banner{}, errGroup.NotFound
	}
	if err != nil {
		return domain.Banner{}, err
	}

	return banner, nil
}

func (s *BannerService) GetAllBannersService(req domain.GetBannersRequest) ([]domain.Banner, error) {
	banners, err := s.BannerRepository.GetAllBannersRepo(req)

	log.Println(banners, err)
	if err == sql.ErrNoRows {
		return nil, errGroup.NotFound
	}
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (s *BannerService) PostBannerService(req domain.Banner) (domain.Banner, error) {
	log.Println(req)
	_, err := s.BannerRepository.GetUserBannerRepo(domain.GetUserBannerRequest{Banner: req, UseLastRevision: true})
	log.Println(err)
	if err != nil && err != sql.ErrNoRows {
		return domain.Banner{}, errGroup.AlreadyExists
	}

	banner, err := s.BannerRepository.PostBannerRepo(req)
	if err != nil {
		return domain.Banner{}, err
	}

	return banner, nil
}

func (s *BannerService) PatchBannerService(req domain.Banner) error {
	err := s.BannerRepository.PatchBannerRepo(req)
	if err == sql.ErrNoRows {
		return errGroup.NotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *BannerService) DeleteBannerService(id int) error {
	err := s.BannerRepository.DeleteBannerRepo(id)
	if err == sql.ErrNoRows {
		return errGroup.NotFound
	}
	if err != nil {
		return err
	}

	return nil
}
