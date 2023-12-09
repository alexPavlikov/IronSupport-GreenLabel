package app

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/chat/internal/entity/message"
	message_db "github.com/alexPavlikov/IronSupport-GreenLabel/chat/internal/entity/message/db"
	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/server"

	"github.com/julienschmidt/httprouter"
)

//var ClientPostgreSQL dbClient.Client

func Run(router *httprouter.Router) *httprouter.Router {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")

	// cfg := config.GetConfig()

	//var err error

	// ClientPostgreSQL, err = dbClient.NewClient(context.TODO(), cfg.Storage)
	// if err != nil {
	// 	logger.Fatalf("failed to get new client postgresql, due to err: %v", err)
	// }

	mRep := message_db.NewRepository(server.ClientPostgreSQL, logger)
	mSer := message.NewService(mRep, logger)
	mHan := message.NewHandler(mSer, logger)

	mHan.Register(router)

	return router
}
