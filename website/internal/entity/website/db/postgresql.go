package site_db

import (
	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) site.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}
