package service

import (
	"context"
	"errors"
	"github.com/Devskare/Warehouse/mocks"
	"github.com/Devskare/Warehouse/modules/Whouse/models"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorageADD_InvalidWeight(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewStorageService(repo, slog.Default())

	err := service.StorageADD(context.Background(), -1)

	require.Error(t, err)
	require.EqualError(t, err, "MaxWeight must be >= 0")

	repo.AssertNotCalled(t, "StorageADD", mock.Anything, mock.Anything)
}

func TestStorageADD_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewStorageService(repo, slog.Default())

	repo.On("StorageADD", mock.Anything, 100.0).
		Return(errors.New("db error"))

	err := service.StorageADD(context.Background(), 100.0)

	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "StorageADD", mock.Anything, 100.0)
}

func TestStorageADD_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewStorageService(repo, slog.Default())

	repo.On("StorageADD", mock.Anything, 100.0).
		Return(nil)

	err := service.StorageADD(context.Background(), 100.0)

	require.NoError(t, err)

	repo.AssertCalled(t, "StorageADD", mock.Anything, 100.0)
}

func TestListStorages_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewStorageService(repo, slog.Default())

	repo.On("ListStorages", mock.Anything).
		Return(nil, errors.New("db error"))

	result, err := service.ListStorages(context.Background())

	require.Error(t, err)
	require.EqualError(t, err, "db error")
	require.Nil(t, result)

	repo.AssertCalled(t, "ListStorages", mock.Anything)
}

func TestListStorages_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewStorageService(repo, slog.Default())

	expected := []models.StorageModel{
		{ID: 1, MaxWeight: 100},
		{ID: 2, MaxWeight: 200},
	}

	repo.On("ListStorages", mock.Anything).
		Return(expected, nil)

	result, err := service.ListStorages(context.Background())

	require.NoError(t, err)
	require.Equal(t, expected, result)

	repo.AssertCalled(t, "ListStorages", mock.Anything)
}
