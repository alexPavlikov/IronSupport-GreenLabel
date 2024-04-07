package news

import (
	"context"
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/news", h.NewsHandler)
	router.HandlerFunc(http.MethodGet, "/post", h.PostHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) NewsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website NewsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	news, err := h.service.GetNews(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"News": news}
	title := map[string]interface{}{"Guest": guest.Guest, "Title": "Новости"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "news", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website NewsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website PostHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	title := map[string]interface{}{"Guest": guest.Guest, "Title": "Публикация №1"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "post", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website PostHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}
