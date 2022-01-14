package nats

import (
	"github.com/TheTeaParty/notnotes-platform/internal/config"
	"github.com/nats-io/nats.go"
)

func New(c *config.Config) (*nats.Conn, error) {
	nc, err := nats.Connect(c.NATSHost)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
