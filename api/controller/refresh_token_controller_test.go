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
)

func TestRefreshTokenController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access-secret",
		RefreshTokenSecret:     "refresh-secret",
		AccessTokenExpiryHour:  1,
		RefreshTokenExpiryHour: 1,
	}

	userID := primitive.NewObjectID().Hex()
	user := domain.User{ID: primitive.NewObjectID(), Name: "Test", Email: "test@example.com"}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		mockUsecase.On("ExtractIDFromToken", "refresh", env.RefreshTokenSecret).Return(userID, nil).Once()
		mockUsecase.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
		mockUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("refresh2", nil).Once()

		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader("refreshToken=refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp domain.RefreshTokenResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "access", resp.AccessToken)
		assert.Equal(t, "refresh2", resp.RefreshToken)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("extract id error", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		mockUsecase.On("ExtractIDFromToken", "refresh", env.RefreshTokenSecret).Return("", errors.New("invalid")).Once()

		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader("refreshToken=refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("get user error", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		mockUsecase.On("ExtractIDFromToken", "refresh", env.RefreshTokenSecret).Return(userID, nil).Once()
		mockUsecase.On("GetUserByID", mock.Anything, userID).Return(domain.User{}, errors.New("not found")).Once()

		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader("refreshToken=refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("create access token error", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		mockUsecase.On("ExtractIDFromToken", "refresh", env.RefreshTokenSecret).Return(userID, nil).Once()
		mockUsecase.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
		mockUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("", errors.New("token err")).Once()

		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader("refreshToken=refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("create refresh token error", func(t *testing.T) {
		mockUsecase := new(mocks.RefreshTokenUsecase)
		mockUsecase.On("ExtractIDFromToken", "refresh", env.RefreshTokenSecret).Return(userID, nil).Once()
		mockUsecase.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
		mockUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("", errors.New("token err")).Once()

		r := gin.New()
		rtc := &controller.RefreshTokenController{RefreshTokenUsecase: mockUsecase, Env: env}
		r.POST("/refresh", rtc.RefreshToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader("refreshToken=refresh"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}
