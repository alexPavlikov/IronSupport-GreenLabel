package site

import (
	"fmt"
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
	//router.HandlerFunc(http.MethodGet, "/account", h.AccountHandler)
	//router.HandlerFunc(http.MethodGet, "/account/notifications", h.NotificationsHandler)
	// router.HandlerFunc(http.MethodGet, "/news", h.NewsHandler)
	// router.HandlerFunc(http.MethodGet, "/post", h.PostHandler)
	// router.HandlerFunc(http.MethodGet, "/products", h.ProductsHandler)
	// router.HandlerFunc(http.MethodGet, "/products/sort", h.ProductsSortHandler)
	//router.HandlerFunc(http.MethodGet, "/clients", h.ClientsHandler)
	router.HandlerFunc(http.MethodGet, "/vacancy", h.VacancyHandler)
	router.HandlerFunc(http.MethodGet, "/about", h.AboutHandler)
	// router.HandlerFunc(http.MethodGet, "/backet", h.BacketHandler)
	// router.HandlerFunc(http.MethodGet, "/purchases", h.PurchasesHandler)

	router.HandlerFunc(http.MethodPost, "/find", h.FindHandler)
	router.HandlerFunc(http.MethodGet, "/find", h.FindResultHandler)

	router.HandlerFunc(http.MethodPost, "/subscribe", h.SubHandler)

	// router.HandlerFunc(http.MethodGet, "/website/menu", h.WebsiteMenuHandler)
	// router.HandlerFunc(http.MethodGet, "/administrator", h.AdminHandler)

	// router.HandlerFunc(http.MethodGet, "/exit", h.ExitHandler)
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

func (h *handler) VacancyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website VacancyHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "vacancy", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website VacancyHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AboutHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "about", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open website AboutHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) FindHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	find := r.FormValue("find")
	if find != "" {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", find)

		http.Redirect(w, r, "/find", http.StatusSeeOther)
	}
}

func (h *handler) FindResultHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "find", nil)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) SubHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mail := r.FormValue("email")
	if mail != "" {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", mail)
		http.Redirect(w, r, "/#email", http.StatusSeeOther)
	}
}
