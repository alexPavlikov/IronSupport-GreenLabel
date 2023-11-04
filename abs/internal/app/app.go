package app

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/abs/internal/entity/post"
	post_db "github.com/alexPavlikov/IronSupport-GreenLabel/abs/internal/entity/post/db"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

var ClientPostgreSQL dbClient.Client

func Run(router *httprouter.Router) *httprouter.Router {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")

	cfg := config.GetConfig()

	var err error

	ClientPostgreSQL, err = dbClient.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("failed to get new client postgresql, due to err: %v", err)
	}

	pRep := post_db.NewRepository(ClientPostgreSQL, logger)
	pSer := post.NewService(pRep, logger)
	pHan := post.NewHandler(pSer, logger)

	pHan.Register(router)

	return router
}
