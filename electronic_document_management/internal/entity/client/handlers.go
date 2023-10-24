package client

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/client", h.ClientHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ClientHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	clients, err := h.service.GetClients(context.TODO())
	if err != nil {
		fmt.Println("Errro", err)
	}

	title := map[string]interface{}{"Title": "ЭДО - Клиенты"}
	data := map[string]interface{}{"Clients": clients}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "client", data)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}
