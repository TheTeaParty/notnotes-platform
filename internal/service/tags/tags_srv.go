package tags

import (
	"context"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"github.com/TheTeaParty/notnotes-platform/internal/service"
	"go.uber.org/zap"
)

type tagsSrv struct {
	tagsRepository domain.TagRepository
	l              *zap.Logger
}

func (s *tagsSrv) GetAll(ctx context.Context, name string) ([]*domain.Tag, error) {
	return s.tagsRepository.GetAll(ctx, name)
}

func NewSrv(tagsRepository domain.TagRepository,
	l *zap.Logger) service.Tags {

	return &tagsSrv{
		tagsRepository: tagsRepository,
		l:              l,
	}
}
