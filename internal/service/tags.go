package service

import (
	"context"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
)

type Tags interface {
	GetAll(ctx context.Context, name string) ([]*domain.Tag, error)
}
