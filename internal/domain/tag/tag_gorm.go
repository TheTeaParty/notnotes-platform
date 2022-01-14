package tag

import (
	"context"
	"fmt"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"gorm.io/gorm"
)

type tagGorm struct {
	db *gorm.DB
}

func (r *tagGorm) Get(ctx context.Context, ID string) (*domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (r *tagGorm) GetByNames(ctx context.Context, names []string) ([]*domain.Tag, error) {

	db := r.db.WithContext(ctx)

	var tags []*domain.Tag
	if err := db.Where("name in ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagGorm) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (r *tagGorm) GetAll(ctx context.Context, name string) ([]*domain.Tag, error) {

	db := r.db.WithContext(ctx)

	if name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%v%%", name))
	}

	var tags []*domain.Tag

	if err := db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagGorm) Create(ctx context.Context, tag *domain.Tag) error {
	//TODO implement me
	panic("implement me")
}

func (r *tagGorm) Update(ctx context.Context, ID string, tag *domain.Tag) error {
	//TODO implement me
	panic("implement me")
}

func (r *tagGorm) Delete(ctx context.Context, ID string) error {
	//TODO implement me
	panic("implement me")
}

func NewGorm(db *gorm.DB) domain.TagRepository {
	return &tagGorm{
		db: db,
	}
}
