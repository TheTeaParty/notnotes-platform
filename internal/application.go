package internal

import (
	"net/http"

	"github.com/TheTeaParty/notnotes-platform/internal/config"
	"github.com/TheTeaParty/notnotes-platform/internal/handler/ws"
	v1 "github.com/TheTeaParty/notnotes-platform/pkg/api/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Application struct {
	r  http.Handler
	ws *ws.WSHandler
	c  *config.Config
	l  *zap.Logger
}

func NewApplication(restHandler v1.ServerInterface, c *config.Config, l *zap.Logger, ws *ws.WSHandler) (*Application, error) {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)

	swagger, _ := v1.GetSwagger()
	swagger.Servers = nil

	h := v1.HandlerFromMux(restHandler, r)
	r.Mount("/", h)

	r.HandleFunc("/ws", ws.ServeWs)

	return &Application{
		r:  r,
		c:  c,
		l:  l,
		ws: ws,
	}, nil
}

func (a *Application) RunHTTP() error {
	if err := http.ListenAndServe(a.c.Port, a.r); err != nil {
		return err
	}

	return nil
}
