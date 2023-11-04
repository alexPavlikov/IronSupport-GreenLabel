package mainABS

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/abs/internal/app"
	"github.com/julienschmidt/httprouter"
)

func MainABS(router *httprouter.Router) *httprouter.Router {
	router = app.Run(router)
	return router
}
