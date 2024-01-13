package client

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

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/client", h.ClientHandler)
	router.HandlerFunc(http.MethodGet, "/edm/client/add", h.AddClientHandler)

	router.HandlerFunc(http.MethodPost, "/edm/client/sorted/", h.SorterClientHandler)
	router.HandlerFunc(http.MethodGet, "/edm/client/sorted/", h.SorterClientHandler)

	router.HandlerFunc(http.MethodGet, "/edm/client/edit", h.EditClientHandler)
	router.HandlerFunc(http.MethodGet, "/edm/client/edits", h.EditPostClientHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ClientHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open ClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	clients, err := h.service.GetClients(context.TODO())
	if err != nil {
		fmt.Println("Errro", err)
	}

	title := map[string]interface{}{"Title": "ЭДО - Клиенты", "Page": "Client"}
	data := map[string]interface{}{"Clients": clients}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open ClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "client", data)
	if err != nil {
		h.logger.Tracef("%s - failed open ClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) AddClientHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var client Client

	client.INN = r.FormValue("inn")
	client.KPP = r.FormValue("kpp")
	client.OGRN = r.FormValue("ogrn")
	client.Name = r.FormValue("name")
	client.Owner = r.FormValue("owner")
	client.Phone = r.FormValue("phone")
	client.Email = r.FormValue("email")
	client.Address = r.FormValue("address")

	err := h.service.AddClient(context.TODO(), &client)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	http.Redirect(w, r, "/edm/client", http.StatusSeeOther)
}

var clients []Client

func (h *handler) SorterClientHandler(w http.ResponseWriter, r *http.Request) {
	var cl Client

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		cl.Name = r.FormValue("name")
		cl.INN = r.FormValue("inn")
		cl.OGRN = r.FormValue("ogrn")
		fmt.Println(cl)

		clients, err = h.service.GetClientsBySorted(context.TODO(), cl)
		if err != nil {
			http.NotFound(w, r)
		}
		fmt.Println(clients)

		http.Redirect(w, r, "/edm/client/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SortedClientHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		fmt.Println(clients)

		data := map[string]interface{}{"Clients": clients}
		header := map[string]string{"Title": "ЭДО - Клиенты", "Page": "Client"}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open SortedClientHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "client", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SortedClientHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		// err = tmpl.ExecuteTemplate(w, "dialog", dialog)
		// if err != nil {
		// 	h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		// 	//http.NotFound(w, r)
		// }
	}

}

func (h *handler) EditClientHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	fmt.Println(idx)

	client, err := h.service.GetClient(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open EditClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	fmt.Println(client)

	title := map[string]string{"Title": "ЭДО - Редактирование клиента", "Page": "Client"}
	data := map[string]interface{}{"Client": client}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "client_edit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditClientHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostClientHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var client Client

	client.Id, _ = strconv.Atoi(r.FormValue("id"))
	status := r.FormValue("status")
	if status == "true" {
		client.Status = true
	} else {
		client.Status = false
	}
	client.INN = r.FormValue("inn")
	client.KPP = r.FormValue("kpp")
	client.OGRN = r.FormValue("ogrn")
	client.Name = r.FormValue("name")
	client.Owner = r.FormValue("owner")
	client.Phone = r.FormValue("phone")
	client.Email = r.FormValue("email")
	client.Address = r.FormValue("address")

	err := h.service.UpdateClient(context.TODO(), &client)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/client", http.StatusSeeOther)
}
