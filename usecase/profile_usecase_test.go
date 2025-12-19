package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain/mocks"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestProfileUsecase_GetProfileByID(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	userObjectID := bson.NewObjectID()
	userID := userObjectID.Hex()

	t.Run("success", func(t *testing.T) {
		mockUser := domain.User{ID: userObjectID, Name: "Test", Email: "test@example.com"}
		mockUserRepository.On("GetByID", mock.Anything, userID).Return(mockUser, nil).Once()

		u := usecase.NewProfileUsecase(mockUserRepository, time.Second*2)
		profile, err := u.GetProfileByID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, mockUser.Name, profile.Name)
		assert.Equal(t, mockUser.Email, profile.Email)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepository.On("GetByID", mock.Anything, userID).Return(domain.User{}, errors.New("Unexpected")).Once()

		u := usecase.NewProfileUsecase(mockUserRepository, time.Second*2)
		profile, err := u.GetProfileByID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, profile)

		mockUserRepository.AssertExpectations(t)
	})
}
