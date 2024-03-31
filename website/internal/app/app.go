package website_app

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/server"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/admin"
	admin_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/admin/db"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	guest_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest/db"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
	news_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news/db"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	product_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product/db"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
	site_db "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website/db"
	"github.com/julienschmidt/httprouter"
)

func Run(router *httprouter.Router) *httprouter.Router {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")

	logger.Info(config.LOG_INFO, " - Start website handlers")
	wRep := site_db.NewRepository(server.ClientPostgreSQL, logger)
	wSer := site.NewService(wRep, logger)
	wHan := site.NewHandler(wSer, logger)
	wHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start website product handlers")
	pRep := product_db.NewRepository(server.ClientPostgreSQL, logger)
	pSer := product.NewService(pRep, logger)
	pHan := product.NewHandler(pSer, logger)
	pHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start website admin handlers")
	aRep := admin_db.NewRepository(server.ClientPostgreSQL, logger)
	aSer := admin.NewService(aRep, logger)
	aHan := admin.NewHandler(aSer, logger)
	aHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start website news handlers")
	nRep := news_db.NewRepository(server.ClientPostgreSQL, logger)
	nSer := news.NewService(nRep, logger)
	nHan := news.NewHandler(nSer, logger)
	nHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start website guest handlers")
	gRep := guest_db.NewRepository(server.ClientPostgreSQL, logger)
	gSer := guest.NewService(gRep, logger)
	gHan := guest.NewHandler(gSer, logger)
	gHan.Register(router)

	return router
}
