package domain

import (
	"context"
	"time"
)

type Tag struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type TagRepository interface {
	Get(ctx context.Context, ID string) (*Tag, error)
	GetByName(ctx context.Context, name string) (*Tag, error)
	GetAll(ctx context.Context, name string) ([]*Tag, error)
	GetByNames(ctx context.Context, names []string) ([]*Tag, error)
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, ID string, tag *Tag) error
	Delete(ctx context.Context, ID string) error
}
