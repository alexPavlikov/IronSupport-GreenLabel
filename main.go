package main

import (
	mainABS "github.com/alexPavlikov/IronSupport-GreenLabel/abs/cmd/web"
	mainChat "github.com/alexPavlikov/IronSupport-GreenLabel/chat/cmd/web"
	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	mainEDO "github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/cmd/web"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/server"
	mainWebsite "github.com/alexPavlikov/IronSupport-GreenLabel/website/cmd/web"
	yandexdisk "github.com/alexPavlikov/IronSupport-GreenLabel/yandex_disk"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	cfg := config.GetConfig()

	router = mainEDO.MainEDO(router)
	router = mainChat.MainChat(router)
	router = mainABS.MainABS(router)
	router = mainWebsite.MainWebsite(router)

	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")
	logger.Info(config.LOG_INFO, " - Start requests handlers")
	isglHan := server.NewHandler(logger)
	isglHan.Register(router)

	yandexdisk.MainYDisk()

	server.Start(router, *cfg)
}
