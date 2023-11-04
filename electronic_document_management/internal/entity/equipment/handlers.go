package equipment

import (
	"context"
	"fmt"
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

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/equipment", h.EquipmentHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) EquipmentHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	eqs, err := h.service.GetEquipments(context.TODO())
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}

	title := map[string]string{"Title": "ЭДО - Оборудование"}
	data := map[string]interface{}{"Equipments": eqs}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "equipment", data)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
}
