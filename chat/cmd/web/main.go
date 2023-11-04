package mainChat

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/chat/internal/app"
	"github.com/julienschmidt/httprouter"
)

func MainChat(router *httprouter.Router) *httprouter.Router {
	router = app.Run(router)
	return router
}
