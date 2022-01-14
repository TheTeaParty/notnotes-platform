package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type Note struct {
	ID        string
	Name      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []*Tag `gorm:"many2many:note_tags;"`
}

type NoteCriteria struct {
	Name     string
	TagNames []string
}

type NoteRepository interface {
	Get(ctx context.Context, ID string) (*Note, error)
	GetByCriteria(ctx context.Context, criteria NoteCriteria) ([]*Note, error)
	Create(ctx context.Context, note *Note, tagNames []string) error
	Update(ctx context.Context, ID string, note *Note, tagNames []string) error
	Delete(ctx context.Context, ID string) error
}
