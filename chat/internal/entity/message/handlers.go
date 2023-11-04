package message

import (
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	router.HandlerFunc(http.MethodGet, "/chat", h.ChatHandler)

}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./chat/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "chat", nil)
	if err != nil {
		http.NotFound(w, r)
	}
}
