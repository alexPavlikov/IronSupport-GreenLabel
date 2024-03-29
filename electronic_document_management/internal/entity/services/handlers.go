package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/service", h.ServicesHandler)
	router.HandlerFunc(http.MethodGet, "/edm/service/add", h.AddServiceHandler)
	router.HandlerFunc(http.MethodGet, "/edm/service/type/add", h.AddTypeServiceHandler)
	router.HandlerFunc(http.MethodPost, "/edm/service/sorted", h.SortServiceHandler)
	router.HandlerFunc(http.MethodGet, "/edm/service/sorted", h.SortServiceHandler)
	router.HandlerFunc(http.MethodGet, "/edm/service/edit", h.EditServiceHandler)
	router.HandlerFunc(http.MethodGet, "/edm/service/edits", h.EditPostServiceHandler)

	router.HandlerFunc(http.MethodGet, "/edm/service/find", h.ServicesFindHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var Events []string

func (h *handler) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	if !user.UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		services, err := h.service.GetServices(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		types, err := h.service.GetServiceType(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		eq, err := h.service.GetAllEquipment(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		Events, err = utils.ReadEventFile()
		if err != nil {
			fmt.Println(err)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Услуги", "Page": "Service", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"Services": services, "Type": types, "Eq": eq, "OK": false, "Auth": arr[2]}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "service", data)
		if err != nil {
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}
}

func (h *handler) AddServiceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var s Services

	s.Equipment, _ = strconv.Atoi(r.FormValue("equipment"))
	s.Type = r.FormValue("type")
	s.Cost, _ = strconv.Atoi(r.FormValue("cost"))

	err := h.service.AddServices(context.TODO(), &s)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/service", http.StatusSeeOther)
}

func (h *handler) AddTypeServiceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")

	fmt.Println(name)

	err := h.service.AddServiceType(context.TODO(), name)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/service", http.StatusSeeOther)
}

var services []Services

func (h *handler) SortServiceHandler(w http.ResponseWriter, r *http.Request) {
	var srv Services

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		srv.Equipment, _ = strconv.Atoi(r.FormValue("equipment"))
		srv.Type = r.FormValue("type")

		fmt.Println(srv)

		services, err = h.service.GetServiceBySort(context.TODO(), &srv)
		if err != nil {
			http.NotFound(w, r)
		}
		fmt.Println(services)

		http.Redirect(w, r, "/edm/service/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SortServiceHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		fmt.Println(services)

		types, err := h.service.GetServiceType(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		eq, err := h.service.GetAllEquipment(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Услуги", "Page": "Service", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"Services": services, "Type": types, "Eq": eq, "Auth": arr[2]}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			h.logger.Tracef("%s - failed open SortServiceHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "service", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SortServiceHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}
	}
}

func (h *handler) EditServiceHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditServiceHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	fmt.Println(idx)

	services, err := h.service.GetService(context.TODO(), idx)
	if err != nil {
		http.NotFound(w, r)
	}

	types, err := h.service.GetServiceType(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	eq, err := h.service.GetAllEquipment(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Услуги", "Page": "Service", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Ser": services, "Type": types, "Eq": eq, "OK": true, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditServiceHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "service_edit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditServiceHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostServiceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var ser Services

	ser.Id, _ = strconv.Atoi(r.FormValue("id"))
	ser.Equipment, _ = strconv.Atoi(r.FormValue("equipment"))
	ser.Type = r.FormValue("type")
	ser.Cost, _ = strconv.Atoi(r.FormValue("cost"))

	fmt.Println(ser)

	err := h.service.UpdateServices(context.TODO(), &ser)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/service", http.StatusSeeOther)
}

func (h *handler) ServicesFindHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	text := r.FormValue("text")

	sr, err := h.service.FindService(context.TODO(), text)
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Поиск", "Page": "Service", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Text": text, "Cat": "Service", "Services": sr, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "find", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
