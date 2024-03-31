package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	router.HandlerFunc(http.MethodGet, "/website/menu", h.WebsiteMenuHandler)
	router.HandlerFunc(http.MethodGet, "/administrator", h.AdminHandler)
	router.HandlerFunc(http.MethodGet, "/administrator/news/edit", h.AdminNewsEdit)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
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

	r.ParseForm()

	page := r.FormValue("page")
	if page == "" {
		page = "main"
	}

	news, err := h.service.GetNews(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"Body": page, "News": news}

	err = tmpl.ExecuteTemplate(w, "panel", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) AdminNewsEdit(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	fmt.Println(id)

	news, err := h.service.GetNewsById(context.TODO(), id)
	if err != nil {
		http.Redirect(w, r, "/administrator", http.StatusSeeOther)
	}

	data := map[string]interface{}{"News": news}

	err = tmpl.ExecuteTemplate(w, "newsedit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

// ip, _, err := net.SplitHostPort(r.RemoteAddr)
// if err != nil {
// 	fmt.Fprintf(w, "ip: %q is not IP:port", r.RemoteAddr)
// 	return
// }

// fmt.Println(r.RemoteAddr, ip)

// if ip == `IP который вам нуежен` {
// 	http.Redirect(w, r, `URL 1`, http.StatusFound)
// } else {
// 	http.Redirect(w, r, `URL 2`, http.StatusFound)
// }
