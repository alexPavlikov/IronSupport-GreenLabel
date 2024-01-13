package objects

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

	router.HandlerFunc(http.MethodGet, "/edm/object", h.ObjectHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/add", h.AddObjectHandler)

	router.HandlerFunc(http.MethodPost, "/edm/object/sorted/", h.SorterObjectHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/sorted/", h.SorterObjectHandler)

	router.HandlerFunc(http.MethodGet, "/edm/object/edit", h.ObjectEditHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/edits", h.ObjectEditsHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal//html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	objs, err := h.service.GetObjects(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	clt, err := h.service.GetClient(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	title := map[string]string{"Title": "ЭДО - Объекты", "Page": "Object"}
	data := map[string]interface{}{"Objs": objs, "Clients": clt}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "object", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) AddObjectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var object Object

	object.Name = r.FormValue("name")
	object.Address = r.FormValue("address")
	object.WorkSchedule = r.FormValue("time")
	object.Client.Id, _ = strconv.Atoi(r.FormValue("client"))

	err := h.service.AddObject(context.TODO(), &object)
	if err != nil {
		http.NotFound(w, r)
	}
	http.Redirect(w, r, "/edm/object", http.StatusSeeOther)
}

var objects []Object

func (h *handler) SorterObjectHandler(w http.ResponseWriter, r *http.Request) {
	var ob Object

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		ob.Name = r.FormValue("name")
		ob.Client.Id, _ = strconv.Atoi(r.FormValue("client"))

		objects, err = h.service.GetObjectBySorted(context.TODO(), &ob)
		if err != nil {
			http.NotFound(w, r)
		}

		http.Redirect(w, r, "/edm/object/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SorterObjectHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		clt, err := h.service.GetClient(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		data := map[string]interface{}{"Objs": objects, "Clients": clt}
		header := map[string]string{"Title": "ЭДО - Объекты", "Page": "Object"}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open SorterObjectHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "object", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SorterObjectHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		// err = tmpl.ExecuteTemplate(w, "dialog", dialog)
		// if err != nil {
		// 	h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		// 	//http.NotFound(w, r)
		// }
	}
}

func (h *handler) ObjectEditHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open ObjectEditHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()

	idx, _ := strconv.Atoi(r.FormValue("id"))

	fmt.Println(idx)

	editobj, err := h.service.GetObject(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open ObjectEditHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	clt, err := h.service.GetClient(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Println(editobj)

	data := map[string]interface{}{"Objectedit": editobj, "Clients": clt}
	header := map[string]string{"Title": "ЭДО - Редактирование объекта", "Page": "Object"}

	err = tmpl.ExecuteTemplate(w, "header", header)
	if err != nil {
		h.logger.Tracef("%s - failed open ObjectEditHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "objectedit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open ObjectEditHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) ObjectEditsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var obj Object

	obj.Id, _ = strconv.Atoi(r.FormValue("id"))
	obj.Name = r.FormValue("name")
	obj.Address = r.FormValue("address")
	obj.WorkSchedule = r.FormValue("time")
	obj.Client.Id, _ = strconv.Atoi(r.FormValue("client"))

	err := h.service.UpdateObject(context.TODO(), &obj)
	if err != nil {
		h.logger.Tracef("%s - failed open ObjectEditsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}
