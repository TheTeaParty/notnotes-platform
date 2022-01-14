package service

import (
	"context"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
)

type Notes interface {
	Create(ctx context.Context, note *domain.Note) error
	Update(ctx context.Context, ID string, note *domain.Note) error
	Delete(ctx context.Context, ID string) error
	Get(ctx context.Context, ID string) (*domain.Note, error)
	GetMatching(ctx context.Context, name string, tagNames []string) ([]*domain.Note, error)
}
