package notes

import (
	"context"
	"regexp"

	"go.uber.org/zap"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"github.com/TheTeaParty/notnotes-platform/internal/service"
)

type notesSrv struct {
	coordinatorService service.Coordinator

	notesRepository domain.NoteRepository
	tagsRepository  domain.TagRepository

	l *zap.Logger
}

func (s *notesSrv) Create(ctx context.Context, note *domain.Note) error {

	tagNames := s.parseTagNames(ctx, note.Content)

	if err := s.notesRepository.Create(ctx, note, tagNames); err != nil {
		return err
	}

	return nil
}

func (s *notesSrv) Update(ctx context.Context, ID string, note *domain.Note) error {
	tagNames := s.parseTagNames(ctx, note.Content)

	if err := s.notesRepository.Update(ctx, ID, note, tagNames); err != nil {
		return err
	}

	return nil
}

func (s *notesSrv) parseTagNames(ctx context.Context, content string) []string {
	var re = regexp.MustCompile(`(?m)\#([\w\n]+)`)

	submatches := re.FindAllStringSubmatch(content, -1)

	matches := make([]string, 0)
	for _, m := range submatches {
		matches = append(matches, m[1])
	}

	return matches
}

func (s *notesSrv) Delete(ctx context.Context, ID string) error {
	err := s.notesRepository.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *notesSrv) Get(ctx context.Context, ID string) (*domain.Note, error) {
	return s.notesRepository.Get(ctx, ID)
}

func (s *notesSrv) GetMatching(ctx context.Context, name string, tagNames []string) ([]*domain.Note, error) {
	return s.notesRepository.GetByCriteria(ctx, domain.NoteCriteria{
		Name:     name,
		TagNames: tagNames,
	})
}

func NewSrv(notesRepository domain.NoteRepository,
	tagsRepository domain.TagRepository,
	l *zap.Logger, coordinatorService service.Coordinator) service.Notes {
	return &notesSrv{
		coordinatorService: coordinatorService,

		notesRepository: notesRepository,
		tagsRepository:  tagsRepository,

		l: l,
	}
}
