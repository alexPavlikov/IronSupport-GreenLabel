package website_app

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/server"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
	website_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website/db"
	"github.com/julienschmidt/httprouter"
)

func Run(router *httprouter.Router) *httprouter.Router {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")

	// cfg := config.GetConfig()

	// var err error

	// ClientPostgreSQL, err = dbClient.NewClient(context.TODO(), cfg.Storage)
	// if err != nil {
	// 	logger.Fatalf("failed to get new client postgresql, due to err: %v", err)
	// }

	logger.Info(config.LOG_INFO, " - Start website handlers")
	wRep := website_db.NewRepository(server.ClientPostgreSQL, logger)
	wSer := website.NewService(wRep, logger)
	wHan := website.NewHandler(wSer, logger)
	wHan.Register(router)

	return router
}
