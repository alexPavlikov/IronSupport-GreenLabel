package guest

import (
	"context"
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
	router.HandlerFunc(http.MethodGet, "/auth", h.AuthHandler)
	router.HandlerFunc(http.MethodPost, "/auth/authconfirm", h.AuthConfirmHandler)

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

var Err bool
var Guest Guests

func (h *handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()
	v := r.FormValue("val")

	text := "singup"
	if v != "" {
		text = "reg"
	}

	// title := map[string]interface{}{}
	data := map[string]interface{}{"Content": text, "Title": "Авторизация", "Err": Err}

	err = tmpl.ExecuteTemplate(w, "auth", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) AuthConfirmHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")
	if email != "" && password != "" {
		var err error
		Guest, err = h.service.AuthGuest(context.TODO(), email, password)
		if err != nil {
			Err = true
			fmt.Println(err)
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
		} else {
			Err = false
			Guest.Auth = true
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		Err = true
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
	}
}

func (h *handler) AccountHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AccountHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	title := map[string]interface{}{"Guest": Guest, "Title": "Аккаунт пользователя"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
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

	tc, err := h.service.GetTrustCompany(context.TODO())
	if err != nil {
		h.logger.Tracef("%s - failed open website ClientsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	fmt.Println("#123123", tc)

	title := map[string]interface{}{"Guest": Guest, "Title": "Наши клиенты"}
	data := map[string]interface{}{"TrustCompany": tc}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "clients", data)
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

	title := map[string]interface{}{"Guest": Guest, "Title": "Уведомления"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
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

	title := map[string]interface{}{"Guest": Guest, "Title": "Корзина"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
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

	title := map[string]interface{}{"Guest": Guest, "Title": "Покупки"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
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
