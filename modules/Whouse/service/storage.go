package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/devskar/warehouse/modules/Whouse/models"
	"github.com/devskar/warehouse/modules/Whouse/repository"
)

type StorageService struct {
	repo repository.WHouser
	log  *slog.Logger
}

func NewStorageService(repo repository.WHouser, log *slog.Logger) *StorageService {
	return &StorageService{repo: repo, log: log}
}

func (s *StorageService) StorageADD(ctx context.Context, MaxWeight float64) error {
	if MaxWeight < 0 {
		s.log.Error("MaxWeight must be >= 0")
		return fmt.Errorf("MaxWeight must be >= 0")
	}
	err := s.repo.StorageADD(ctx, MaxWeight)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageService) ListStorages(ctx context.Context) ([]models.StorageModel, error) {
	storages, err := s.repo.ListStorages(ctx)
	if err != nil {
		s.log.Error("Failed to list storages", "error", err)
		return nil, err
	}
	return storages, nil
}
