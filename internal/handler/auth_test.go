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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{AuthService: mockAuthService, BannerService: mockBannerService}

	req := domain.AuthStruct{}
	reqBody, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqBody)

	mockAuthService.On("RegisterService", req).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register", reqReader)
	c.Request.Header.Set("Content-Type", "application/json")

	h.RegisterHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{AuthService: mockAuthService, BannerService: mockBannerService}

	req := domain.AuthStruct{}
	reqBody, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqBody)

	mockAuthService.On("LoginService", req).Return("token", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", reqReader)
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestRequireAuth(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{AuthService: mockAuthService, BannerService: mockBannerService}

	user := domain.AuthStruct{}
	mockAuthService.On("CheckTokenService", "token").Return(user, nil)

	router := gin.Default()
	router.Use(h.RequireAuth)
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test passed")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	//req.Header.Set("Authorization", "token")
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: "token"})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test passed", w.Body.String())
	mockAuthService.AssertExpectations(t)
}
func TestAdminCheck(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	mockBannerService := new(mocks.BannerService)
	h := handler.Handler{AuthService: mockAuthService, BannerService: mockBannerService}

	user := domain.AuthStruct{IsAdmin: true}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/admin", nil)
	c.Set("user", user)

	h.AdminCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
