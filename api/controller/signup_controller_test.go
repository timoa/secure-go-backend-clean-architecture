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
	"golang.org/x/crypto/bcrypt"
)

func TestSignupController_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access-secret",
		RefreshTokenSecret:     "refresh-secret",
		AccessTokenExpiryHour:  1,
		RefreshTokenExpiryHour: 1,
	}

	t.Run("success", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)

		mockSignupUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, errors.New("not found")).Once()

		mockSignupUsecase.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			if u == nil {
				return false
			}
			if u.Email != "test@example.com" || u.Name != "Test" {
				return false
			}
			if u.Password == "password" {
				return false
			}
			return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password")) == nil
		})).Return(nil).Once()

		mockSignupUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockSignupUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("refresh", nil).Once()

		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		body := strings.NewReader("name=Test&email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/signup", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.SignupResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "access", resp.AccessToken)
		assert.Equal(t, "refresh", resp.RefreshToken)

		mockSignupUsecase.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)
		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)
		mockSignupUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, nil).Once()

		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		body := strings.NewReader("name=Test&email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/signup", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockSignupUsecase.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)
		mockSignupUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, errors.New("not found")).Once()
		mockSignupUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("db error")).Once()

		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		body := strings.NewReader("name=Test&email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/signup", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockSignupUsecase.AssertExpectations(t)
	})

	t.Run("create access token error", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)
		mockSignupUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, errors.New("not found")).Once()
		mockSignupUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		mockSignupUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("", errors.New("token error")).Once()

		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		body := strings.NewReader("name=Test&email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/signup", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockSignupUsecase.AssertExpectations(t)
	})

	t.Run("create refresh token error", func(t *testing.T) {
		mockSignupUsecase := new(mocks.SignupUsecase)
		mockSignupUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, errors.New("not found")).Once()
		mockSignupUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		mockSignupUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access", nil).Once()
		mockSignupUsecase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("", errors.New("token error")).Once()

		r := gin.New()
		sc := &controller.SignupController{SignupUsecase: mockSignupUsecase, Env: env}
		r.POST("/signup", sc.Signup)

		rec := httptest.NewRecorder()
		body := strings.NewReader("name=Test&email=test%40example.com&password=password")
		req := httptest.NewRequest(http.MethodPost, "/signup", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockSignupUsecase.AssertExpectations(t)
	})
}
