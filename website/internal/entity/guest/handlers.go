package guest

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
	router.HandlerFunc(http.MethodGet, "/account", h.AccountHandler)
	router.HandlerFunc(http.MethodGet, "/account/notifications", h.NotificationsHandler)
	router.HandlerFunc(http.MethodGet, "/clients", h.ClientsHandler)
	router.HandlerFunc(http.MethodGet, "/backet", h.BacketHandler)
	router.HandlerFunc(http.MethodGet, "/purchases", h.PurchasesHandler)
	router.HandlerFunc(http.MethodGet, "/exit", h.ExitHandler)

}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) AccountHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AccountHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "account", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website AccountHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) ClientsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website ClientsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "clients", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website ClientsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}
func (h *handler) NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website NotificationsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "notifications", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website NotificationsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) BacketHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website BacketHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "backet", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website BacketHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) PurchasesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website PurchasesHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "purchases", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website PurchasesHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) ExitHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
