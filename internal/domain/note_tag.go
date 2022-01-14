package domain

import (
	"time"
)

type NoteTag struct {
	ID        string
	NoteID    string
	Note      *Note
	TagID     string
	Tag       *Tag
	CreatedAt time.Time
}
