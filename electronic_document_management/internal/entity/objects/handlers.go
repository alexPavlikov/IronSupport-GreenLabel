package objects

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

	router.HandlerFunc(http.MethodGet, "/edm/object", h.ObjectHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/add", h.AddObjectHandler)

	router.HandlerFunc(http.MethodPost, "/edm/object/sorted/", h.SorterObjectHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/sorted/", h.SorterObjectHandler)

	router.HandlerFunc(http.MethodGet, "/edm/object/edit", h.ObjectEditHandler)
	router.HandlerFunc(http.MethodGet, "/edm/object/edits", h.ObjectEditsHandler)

	router.HandlerFunc(http.MethodGet, "/edm/object/find", h.RequestFindHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var Events []string

func (h *handler) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	if !user.UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		Events, err = utils.ReadEventFile()
		if err != nil {
			fmt.Println(err)
		}

		objs, err := h.service.GetObjects(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		clt, err := h.service.GetClient(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Объекты", "Page": "Object", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"Objs": objs, "Clients": clt, "OK": false, "Auth": arr[2]}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "object", data)
		if err != nil {
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
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

		arr := utils.ReadCookies(r)

		data := map[string]interface{}{"Objs": objects, "Clients": clt, "OK": true, "Auth": arr[2]}
		header := map[string]interface{}{"Title": "ЭДО - Объекты", "Page": "Object", "Events": Events, "Auth": arr[2]}
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

	arr := utils.ReadCookies(r)

	data := map[string]interface{}{"Objectedit": editobj, "Clients": clt, "Auth": arr[2]}
	header := map[string]interface{}{"Title": "ЭДО - Редактирование объекта", "Page": "Object", "Events": Events, "Auth": arr[2]}

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

func (h *handler) RequestFindHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	text := r.FormValue("text")

	obj, err := h.service.FindObject(context.TODO(), text)
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Поиск", "Page": "Object", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Text": text, "Cat": "Object", "Objs": obj, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "find", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
