package message

import (
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

	router.HandlerFunc(http.MethodGet, "/chat", h.ChatHandler)
	router.HandlerFunc(http.MethodGet, "/chat/add", h.ChatAddHandler)

}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./chat/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	// NOT WORK

	group := r.FormValue("group")

	if group == "" {
		group = "Clear"
	} else {
		group = "NotClear"
		fmt.Println(group)
	}

	data := map[string]interface{}{"Key": group}
	fmt.Println(group)
	err = tmpl.ExecuteTemplate(w, "chat", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) ChatAddHandler(w http.ResponseWriter, r *http.Request) { // array
	r.ParseForm()

	name := r.FormValue("name")
	one := r.FormValue("one")
	two := r.FormValue("two")
	three := r.FormValue("three")
	fmt.Println(name, one, two, three)
}
