package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheTeaParty/notnotes-platform/internal/domain"
	"github.com/TheTeaParty/notnotes-platform/internal/service"
	v1 "github.com/TheTeaParty/notnotes-platform/pkg/api/openapi"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	clients            map[*Client]bool
	broadcast          chan []byte
	register           chan *Client
	unregister         chan *Client
	coordinatorService service.Coordinator

	l *zap.Logger
}

func (h *WSHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	h.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHandler) noteHandler(eventName string) service.GenericHandler {
	return func(ctx context.Context, note *domain.Note) error {

		response := v1.Event{
			Data: v1.Note{
				Content:   note.Content,
				CreatedAt: note.CreatedAt,
				Id:        note.ID,
				Name:      note.Name,
				UpdatedAt: note.UpdatedAt,
			},
			Type: eventName,
		}

		responseJ, err := json.Marshal(response)
		if err != nil {
			return err
		}

		for c := range h.clients {
			c.send <- responseJ
		}

		return nil
	}
}

func (h *WSHandler) noteCreated(ctx context.Context, note *domain.Note) error {

	response := v1.Event{
		Data: v1.Note{
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			Id:        note.ID,
			Name:      note.Name,
			UpdatedAt: note.UpdatedAt,
		},
		Type: string(service.TriggerNoteCreated),
	}

	responseJ, err := json.Marshal(response)
	if err != nil {
		return err
	}

	for c := range h.clients {
		c.send <- responseJ
	}

	return nil
}

func (h *WSHandler) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewWSHandler(l *zap.Logger, coordinatorService service.Coordinator) *WSHandler {
	h := &WSHandler{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		coordinatorService: coordinatorService,
		l:                  l,
	}

	go h.run()

	h.coordinatorService.Handle(service.TriggerNoteCreated, h.noteHandler(string(service.TriggerNoteCreated)))
	h.coordinatorService.Handle(service.TriggerNoteUpdated, h.noteHandler(string(service.TriggerNoteUpdated)))
	h.coordinatorService.Handle(service.TriggerNoteDeleted, h.noteHandler(string(service.TriggerNoteDeleted)))

	return h
}
