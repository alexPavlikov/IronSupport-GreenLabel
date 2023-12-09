package website

import (
	"net/http"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	// router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/", h.IndexHandler)
	router.HandlerFunc(http.MethodGet, "/news", h.NewsHandler)

	router.HandlerFunc(http.MethodGet, "/website/menu", h.WebsiteMenuHandler)
	router.HandlerFunc(http.MethodGet, "/administrator", h.AdminHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "website", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) NewsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website NewsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "news", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website NewsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) WebsiteMenuHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open WebsiteMenuHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "menuweb", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open WebsiteMenuHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "panel", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}
