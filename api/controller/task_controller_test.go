package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/controller"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil).Once()

		r := gin.New()
		r.Use(setUserID(userID))

		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.POST("/task", tc.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader("title=Test"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.SuccessResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Task created successfully", resp.Message)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		r := gin.New()
		r.Use(setUserID(userID))
		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.POST("/task", tc.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		r := gin.New()
		r.Use(setUserID("invalid"))
		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.POST("/task", tc.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader("title=Test"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(errors.New("unexpected")).Once()

		r := gin.New()
		r.Use(setUserID(userID))
		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.POST("/task", tc.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader("title=Test"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestTaskController_Fetch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		mockUsecase.On("FetchByUserID", mock.Anything, userID).Return([]domain.Task{{Title: "T"}}, nil).Once()

		r := gin.New()
		r.Use(setUserID(userID))
		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.GET("/task", tc.Fetch)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/task", nil)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "\"title\"")
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase := new(mocks.TaskUsecase)
		mockUsecase.On("FetchByUserID", mock.Anything, userID).Return(nil, errors.New("unexpected")).Once()

		r := gin.New()
		r.Use(setUserID(userID))
		tc := &controller.TaskController{TaskUsecase: mockUsecase}
		r.GET("/task", tc.Fetch)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/task", nil)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}
