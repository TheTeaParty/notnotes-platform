package rest

import (
	"net/http"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"github.com/TheTeaParty/notnotes-platform/internal/service"
	v1 "github.com/TheTeaParty/notnotes-platform/pkg/api/openapi"
	"github.com/go-chi/render"
)

type restV1Handler struct {
	notesService service.Notes
	tagsService  service.Tags
}

func (h *restV1Handler) GetNotes(w http.ResponseWriter, r *http.Request, params v1.GetNotesParams) {

	var name string
	var tagNames []string

	if params.Name != nil {
		name = *params.Name
	}

	if params.Tags != nil {
		tagNames = *params.Tags
	}

	notes, err := h.notesService.GetMatching(r.Context(), name, tagNames)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	response := make([]*v1.Note, len(notes))
	for i, n := range notes {
		response[i] = &v1.Note{
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			Id:        n.ID,
			Name:      n.Name,
			UpdatedAt: n.UpdatedAt,
		}
	}

	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, response)
}

func (h *restV1Handler) CreateNote(w http.ResponseWriter, r *http.Request) {

	var req v1.CreateNoteJSONBody
	if err := render.DefaultDecoder(r, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	n := &domain.Note{
		Name:    req.Name,
		Content: req.Content,
	}

	if err := h.notesService.Create(r.Context(), n); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.DefaultResponder(w, r, &v1.Note{
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		Id:        n.ID,
		Name:      n.Name,
		UpdatedAt: n.UpdatedAt,
	})
}

func (h *restV1Handler) DeleteNote(w http.ResponseWriter, r *http.Request, noteID string) {
	if err := h.notesService.Delete(r.Context(), noteID); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *restV1Handler) GetNote(w http.ResponseWriter, r *http.Request, noteID string) {
	n, err := h.notesService.Get(r.Context(), noteID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, &v1.Note{
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		Id:        n.ID,
		Name:      n.Name,
		UpdatedAt: n.UpdatedAt,
	})
}

func (h *restV1Handler) UpdateNote(w http.ResponseWriter, r *http.Request, noteID string) {
	var req v1.UpdateNoteJSONBody
	if err := render.DefaultDecoder(r, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	n := &domain.Note{
		Name:    req.Name,
		Content: req.Content,
	}

	if err := h.notesService.Update(r.Context(), noteID, n); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, &v1.Note{
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		Id:        n.ID,
		Name:      n.Name,
		UpdatedAt: n.UpdatedAt,
	})
}

func (h *restV1Handler) GetTags(w http.ResponseWriter, r *http.Request, params v1.GetTagsParams) {

	var name string

	if params.Name != nil {
		name = *params.Name
	}

	tags, err := h.tagsService.GetAll(r.Context(), name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.DefaultResponder(w, r, v1.Error{Message: err.Error()})
		return
	}

	response := make([]*v1.Tag, len(tags))
	for i, t := range tags {
		response[i] = &v1.Tag{
			Id:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
		}
	}

	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, response)
}

func NewRESTV1Handler(
	notesService service.Notes,
	tagsService service.Tags) v1.ServerInterface {
	return &restV1Handler{
		notesService: notesService,
		tagsService:  tagsService,
	}
}
