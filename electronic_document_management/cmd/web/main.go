package mainEDO

import (
	edm_app "github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/app"
	"github.com/julienschmidt/httprouter"
)

func MainEDO(router *httprouter.Router) *httprouter.Router {
	router = edm_app.Run(router)
	return router
}
