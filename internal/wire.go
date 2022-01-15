//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"

	"github.com/TheTeaParty/notnotes-platform/internal/config"
	"github.com/TheTeaParty/notnotes-platform/internal/domain/note"
	"github.com/TheTeaParty/notnotes-platform/internal/domain/tag"
	"github.com/TheTeaParty/notnotes-platform/internal/handler/rest"
	"github.com/TheTeaParty/notnotes-platform/internal/handler/ws"
	"github.com/TheTeaParty/notnotes-platform/internal/pkg/logger"
	"github.com/TheTeaParty/notnotes-platform/internal/pkg/nats"
	"github.com/TheTeaParty/notnotes-platform/internal/pkg/pg"
	"github.com/TheTeaParty/notnotes-platform/internal/service/coordinator"
	"github.com/TheTeaParty/notnotes-platform/internal/service/notes"
	"github.com/TheTeaParty/notnotes-platform/internal/service/tags"
)

func InitializeApplication() (*Application, error) {
	wire.Build(
		config.New,
		logger.New,
		pg.New,
		nats.New,

		note.NewGorm,
		tag.NewGorm,

		notes.NewSrv,
		tags.NewSrv,
		coordinator.NewNATS,

		rest.NewRESTV1Handler,
		ws.NewWSHandler,

		NewApplication,
	)
	return &Application{}, nil
}
