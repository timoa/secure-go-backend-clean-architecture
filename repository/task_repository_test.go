package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo/mocks"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestTaskRepository_Create(t *testing.T) {
	databaseHelper := &mocks.Database{}
	collectionHelper := &mocks.Collection{}

	collectionName := domain.CollectionTask

	mockTask := &domain.Task{ID: bson.NewObjectID(), Title: "Test", UserID: bson.NewObjectID()}

	t.Run("success", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(bson.NewObjectID(), nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		err := r.Create(context.Background(), mockTask)

		assert.NoError(t, err)
		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil, errors.New("Unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		err := r.Create(context.Background(), mockTask)

		assert.Error(t, err)
		collectionHelper.AssertExpectations(t)
	})
}

func TestTaskRepository_FetchByUserID(t *testing.T) {
	databaseHelper := &mocks.Database{}
	collectionHelper := &mocks.Collection{}
	cursorHelper := &mocks.Cursor{}

	collectionName := domain.CollectionTask

	t.Run("invalid user id", func(t *testing.T) {
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		tasks, err := r.FetchByUserID(context.Background(), "invalid")
		assert.Error(t, err)
		assert.Nil(t, tasks)

		databaseHelper.AssertExpectations(t)
	})

	t.Run("find error", func(t *testing.T) {
		userObjectID := bson.NewObjectID()
		userID := userObjectID.Hex()

		collectionHelper.On("Find", mock.Anything, mock.Anything).Return(nil, errors.New("Unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		tasks, err := r.FetchByUserID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, tasks)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		userObjectID := bson.NewObjectID()
		userID := userObjectID.Hex()

		collectionHelper.On("Find", mock.Anything, mock.Anything).Return(cursorHelper, nil).Once()
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			ptr := args.Get(1).(*[]domain.Task)
			*ptr = []domain.Task{{ID: bson.NewObjectID(), Title: "T", UserID: userObjectID}}
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		tasks, err := r.FetchByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, "T", tasks[0].Title)

		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
	})

	t.Run("empty tasks", func(t *testing.T) {
		userObjectID := bson.NewObjectID()
		userID := userObjectID.Hex()

		collectionHelper.On("Find", mock.Anything, mock.Anything).Return(cursorHelper, nil).Once()
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			ptr := args.Get(1).(*[]domain.Task)
			*ptr = nil
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		r := repository.NewTaskRepository(databaseHelper, collectionName)
		tasks, err := r.FetchByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, tasks)
		assert.Len(t, tasks, 0)

		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
	})
}
