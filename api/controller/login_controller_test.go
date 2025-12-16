package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/controller"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access-secret",
		RefreshTokenSecret:     "refresh-secret",
		AccessTokenExpiryHour:  1,
		RefreshTokenExpiryHour: 1,
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := domain.User{ID: primitive.NewObjectID(), Name: "Test", Email: "test@example.com", Password: string(hashed)}

	t.Run("success", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		mockLoginUsecase.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil).Once()
		mockLoginUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockLoginUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("refresh", nil).Once()

		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		body := strings.NewReader("email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.LoginResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "access", resp.AccessToken)
		assert.Equal(t, "refresh", resp.RefreshToken)

		mockLoginUsecase.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		mockLoginUsecase.On("GetUserByEmail", mock.Anything, user.Email).Return(domain.User{}, errors.New("not found")).Once()

		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		body := strings.NewReader("email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockLoginUsecase.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		mockLoginUsecase.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil).Once()

		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		body := strings.NewReader("email=test%40example.com&password=wrong")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockLoginUsecase.AssertExpectations(t)
	})

	t.Run("create access token error", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		mockLoginUsecase.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil).Once()
		mockLoginUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("", errors.New("token error")).Once()

		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		body := strings.NewReader("email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockLoginUsecase.AssertExpectations(t)
	})

	t.Run("create refresh token error", func(t *testing.T) {
		mockLoginUsecase := new(mocks.LoginUsecase)
		mockLoginUsecase.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil).Once()
		mockLoginUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockLoginUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("", errors.New("token error")).Once()

		r := gin.New()
		lc := &controller.LoginController{LoginUsecase: mockLoginUsecase, Env: env}
		r.POST("/login", lc.Login)

		rec := httptest.NewRecorder()
		body := strings.NewReader("email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockLoginUsecase.AssertExpectations(t)
	})
}
