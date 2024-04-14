package handler

import (
	"avito/internal/domain"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	RegisterService(req domain.AuthStruct) error
	LoginService(req domain.AuthStruct) (string, error)
	CheckTokenService(token string) (domain.AuthStruct, error)
}

type BannerService interface {
	GetUserBannerService(req domain.GetUserBannerRequest) (domain.Banner, error)
	GetAllBannersService(req domain.GetBannersRequest) ([]domain.Banner, error)
	PostBannerService(req domain.Banner) (domain.Banner, error)
	PatchBannerService(req domain.Banner) error
	DeleteBannerService(id int) error
}

type Handler struct {
	AuthService   AuthService
	BannerService BannerService
}

func NewHandler(authService AuthService, bannerService BannerService) *gin.Engine {
	g := gin.Default()

	handler := Handler{
		AuthService:   authService,
		BannerService: bannerService,
	}
	g.POST("/login", handler.LoginHandler)
	g.POST("/register", handler.RegisterHandler)

	authApiGroup := g.Group("/").Use(handler.RequireAuth)
	{
		authApiGroup.GET("/user_banner", handler.GetUserBannerHandler)
		authApiGroup.GET("/banner", handler.AdminCheck, handler.GetAllBannersHandler)
		authApiGroup.POST("/banner", handler.AdminCheck, handler.PostBannerHandler)
		authApiGroup.PATCH("/banner/:id", handler.AdminCheck, handler.PatchBannerHandler)
		authApiGroup.DELETE("/banner/:id", handler.AdminCheck, handler.DeleteBannerHandler)
	}

	return g

}
