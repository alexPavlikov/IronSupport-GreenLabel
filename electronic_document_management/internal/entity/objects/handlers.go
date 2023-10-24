package objects

import (
	"context"
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/object", h.ObjectHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	objs, err := h.service.GetObjects(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	title := map[string]string{"Title": "ЭДО - Объекты"}
	data := map[string]interface{}{"Objs": objs}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "object", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
