package mainWebsite

import (
	website_app "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/app"
	"github.com/julienschmidt/httprouter"
)

func MainWebsite(router *httprouter.Router) *httprouter.Router {
	router = website_app.Run(router)
	return router
}
