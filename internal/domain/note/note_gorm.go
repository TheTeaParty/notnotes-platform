package note

import (
	"context"
	"fmt"
	"time"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type noteGorm struct {
	db *gorm.DB
}

func (r *noteGorm) Get(ctx context.Context, ID string) (*domain.Note, error) {

	db := r.db.WithContext(ctx)

	var n *domain.Note

	if err := db.Where("id = ?", ID).First(&n).Error; err != nil {
		return nil, err
	}

	return n, nil
}

func (r *noteGorm) tagNames(names []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(names) > 0 {
			db.Preload("Tags", "name IN ?", names)
		}

		return db.Preload("Tags")
	}
}

func (r *noteGorm) GetByCriteria(ctx context.Context, criteria domain.NoteCriteria) ([]*domain.Note, error) {

	db := r.db.WithContext(ctx)

	if len(criteria.TagNames) > 0 {
		db = db.Preload("Tags", "name IN ?", criteria.TagNames)
	} else {
		db = db.Preload("Tags")
	}

	if criteria.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%v%%", criteria.Name))
	}

	var nts []*domain.Note
	if err := db.Find(&nts).Error; err != nil {
		return nil, err
	}

	ntss := make([]*domain.Note, 0)
	for _, n := range nts {
		if len(n.Tags) > 0 || len(criteria.TagNames) == 0 {
			ntss = append(ntss, n)
		}
	}

	return ntss, nil
}

func (r *noteGorm) Create(ctx context.Context, note *domain.Note, tagNames []string) error {

	db := r.db.WithContext(ctx)

	err := db.Transaction(func(tx *gorm.DB) error {

		note.ID = uuid.New().String()

		if err := tx.Create(&note).Error; err != nil {
			return err
		}

		errors := make(chan error)
		result := make(chan *domain.Tag)

		for _, n := range tagNames {
			go func(n string) {
				tag := &domain.Tag{
					ID:   uuid.New().String(),
					Name: n,
				}

				err := tx.Where("name = ?", n).FirstOrCreate(&tag).Error

				errors <- err
				result <- tag
			}(n)
		}

		for i := 0; i < len(tagNames); i++ {

			err := <-errors
			if err != nil {
				return err
			}

			t := <-result

			nt := &domain.NoteTag{
				ID:     uuid.New().String(),
				TagID:  t.ID,
				NoteID: note.ID,
			}

			err = tx.Create(&nt).Error
			if err != nil {
				return err
			}

		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *noteGorm) Update(ctx context.Context, ID string, note *domain.Note, tagNames []string) error {
	db := r.db.WithContext(ctx)

	err := db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("note_id = ?", ID).Delete(&domain.NoteTag{}).Error; err != nil {
			return err
		}

		updatedAt := time.Now()

		if err := tx.Where("id = ?", ID).Updates(
			&domain.Note{Name: note.Name, Content: note.Content, UpdatedAt: updatedAt},
		).Error; err != nil {
			return err
		}

		note.UpdatedAt = updatedAt
		note.ID = ID

		errors := make(chan error)
		result := make(chan *domain.Tag)

		for _, n := range tagNames {
			go func(n string) {
				tag := &domain.Tag{
					Name: n,
				}

				err := tx.Where(domain.Tag{Name: n}).
					Attrs(domain.Tag{ID: uuid.New().String()}).
					FirstOrCreate(&tag).Error

				errors <- err
				result <- tag
			}(n)
		}

		for i := 0; i < len(tagNames); i++ {

			err := <-errors
			if err != nil {
				return err
			}

			t := <-result

			nt := &domain.NoteTag{
				ID:     uuid.New().String(),
				TagID:  t.ID,
				NoteID: note.ID,
			}

			err = tx.Create(&nt).Error
			if err != nil {
				return err
			}

		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *noteGorm) Delete(ctx context.Context, ID string) error {

	db := r.db.WithContext(ctx)

	db.Where("id = ?", ID)

	if err := db.Where("id = ?", ID).Delete(&domain.Note{}).Error; err != nil {
		return err
	}

	return nil
}

func NewGorm(db *gorm.DB) domain.NoteRepository {
	return &noteGorm{
		db: db,
	}
}
