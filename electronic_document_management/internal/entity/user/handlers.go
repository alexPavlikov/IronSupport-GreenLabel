package user

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

// Register implements handlers.Handlers.
func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/edm/user", h.UserHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	users, err := h.service.GetUsers(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Println(users)

	title := map[string]string{"Title": "ЭДО - Пользователи"}
	data := map[string]interface{}{"User": users}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "user", data)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
}
