package message_db

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/chat/internal/entity/message"
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) message.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}
