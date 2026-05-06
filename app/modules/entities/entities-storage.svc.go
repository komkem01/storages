package entities

import (
	"context"
	"errors"

	"storage/app/modules/entities/ent"
	entitiesinf "storage/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StorageEntity = (*Service)(nil)

func (s *Service) CreateStorage(ctx context.Context, storage *ent.Storage) (*ent.Storage, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}

	_, err := s.db.NewInsert().Model(storage).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Service) GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.Storage, error) {
	storage := new(ent.Storage)
	if err := s.db.NewSelect().Model(storage).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return storage, nil
}
