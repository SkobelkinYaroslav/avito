package handler_test

import (
	"avito/internal/domain"
	"avito/internal/handler"
	mocks "avito/internal/mocks/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserBannerHandler(t *testing.T) {
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{BannerService: mockBannerService}

	banner := domain.Banner{
		TagIDs:    []int{1},
		FeatureID: 2,
		Content:   domain.Content{},
		IsActive:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	req := domain.GetUserBannerRequest{
		Banner:          banner,
		UseLastRevision: false,
	}

	user := domain.AuthStruct{
		ID:       1,
		Email:    "",
		Password: "",
		IsAdmin:  true,
	}

	mockBannerService.On("GetUserBannerService", req).Return(banner, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", user)

	c.Request, _ = http.NewRequest(http.MethodGet, "/user_banner?tag_id=1&feature_id=2&use_last_revision=false", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.GetUserBannerHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBannerService.AssertExpectations(t)
}

func TestGetAllBannersHandler(t *testing.T) {
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{BannerService: mockBannerService}

	req := domain.GetBannersRequest{
		Limit: 10,
	}
	reqBody, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqBody)

	mockBannerService.On("GetAllBannersService", req).Return([]domain.Banner{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/banners", reqReader)
	c.Request.Header.Set("Content-Type", "application/json")

	h.GetAllBannersHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBannerService.AssertExpectations(t)
}

func TestPostBannerHandler(t *testing.T) {
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{BannerService: mockBannerService}

	req := domain.Banner{
		BannerID: 1,
	}
	reqBody, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqBody)

	mockBannerService.On("PostBannerService", req).Return(req, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/banner", reqReader)
	c.Request.Header.Set("Content-Type", "application/json")

	h.PostBannerHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockBannerService.AssertExpectations(t)
}

func TestPatchBannerHandler(t *testing.T) {
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{BannerService: mockBannerService}

	req := domain.Banner{
		BannerID:  1,
		TagIDs:    []int{1},
		FeatureID: 2,
		Content:   domain.Content{},
		IsActive:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	reqBody, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqBody)

	mockBannerService.On("PatchBannerService", req).Return(nil)

	w := httptest.NewRecorder()
	router := gin.Default()
	router.PATCH("/banner/:id", h.PatchBannerHandler)
	request, _ := http.NewRequest(http.MethodPatch, "/banner/1", reqReader)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBannerService.AssertExpectations(t)
}

func TestDeleteBannerHandler(t *testing.T) {
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{BannerService: mockBannerService}

	mockBannerService.On("DeleteBannerService", 1).Return(nil)

	w := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/banner/:id", h.DeleteBannerHandler)
	request, _ := http.NewRequest(http.MethodDelete, "/banner/1", nil)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBannerService.AssertExpectations(t)
}
