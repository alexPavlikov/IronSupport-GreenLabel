package equipment

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
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

	router.HandlerFunc(http.MethodGet, "/edm/equipment", h.EquipmentHandler)
	router.HandlerFunc(http.MethodGet, "/edm/equipment/add", h.AddEquipmentHandler)
	router.HandlerFunc(http.MethodPost, "/edm/equipment/sorted", h.SortHandlerEquipment)
	router.HandlerFunc(http.MethodGet, "/edm/equipment/sorted", h.SortHandlerEquipment)

	router.HandlerFunc(http.MethodGet, "/edm/equipment/edit", h.EditEquipmentHandler)
	router.HandlerFunc(http.MethodGet, "/edm/equipment/edits", h.EditPostEquipmentHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) EquipmentHandler(w http.ResponseWriter, r *http.Request) {
	if !user.UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		eqs, err := h.service.GetEquipments(context.TODO())
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}

		sort, err := h.service.GetAllSortVal(context.TODO())
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}

		title := map[string]string{"Title": "ЭДО - Оборудование", "Page": "Equipment"}
		data := map[string]interface{}{"Equipments": eqs, "Sort": sort, "OK": false}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
		err = tmpl.ExecuteTemplate(w, "equipment", data)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}
}

func (h *handler) AddEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var eq Equipment

	eq.Name = r.FormValue("name")
	eq.Type = r.FormValue("type")
	eq.Manufacture = r.FormValue("manufacture")
	eq.Model = r.FormValue("model")
	eq.UniqueNumber = r.FormValue("number")
	eq.Contract = r.FormValue("file")
	eq.CreateDate = time.Now().Format("02-01-2006")

	fmt.Println(eq.Model)

	err := h.service.AddEquipment(context.TODO(), &eq)
	if err != nil {
		http.NotFound(w, r)
	}
	http.Redirect(w, r, "/edm/equipment", http.StatusSeeOther)
}

var eqs []Equipment

func (h *handler) SortHandlerEquipment(w http.ResponseWriter, r *http.Request) {
	var eq Equipment

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		eq.Name = r.FormValue("name")
		eq.Type = r.FormValue("type")
		eq.Manufacture = r.FormValue("manufacture")
		eq.Model = r.FormValue("model")
		eq.UniqueNumber = r.FormValue("unique")

		if eq.Name != "" {
			eq.Name = "%" + eq.Name + "%"
		}
		if eq.UniqueNumber != "" {
			eq.UniqueNumber = "%" + eq.UniqueNumber + "%"
		}

		fmt.Println(eq)

		eqs, err = h.service.GetEquipmentsBySort(context.TODO(), &eq)
		if err != nil {
			http.NotFound(w, r)
		}
		fmt.Println(eqs)

		http.Redirect(w, r, "/edm/equipment/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SortHandlerEquipment", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		fmt.Println(eqs)

		sort, err := h.service.GetAllSortVal(context.TODO())
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}

		data := map[string]interface{}{"Equipments": eqs, "Sort": sort}
		header := map[string]string{"Title": "ЭДО - Оборудование", "Page": "Equipment"}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open SortHandlerEquipment", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "equipment", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SortHandlerEquipment", config.LOG_ERROR)
			//http.NotFound(w, r)
		}
	}
}

func (h *handler) EditEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditEquipmentHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	fmt.Println(idx)

	eq, err := h.service.GetEquipment(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open EditEquipmentHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	sort, err := h.service.GetAllSortVal(context.TODO())
	if err != nil {
		h.logger.Tracef("%s - failed open EditEquipmentHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	fmt.Println(eq)

	title := map[string]string{"Title": "ЭДО - Редактирование оборудования", "Page": "Equipment"}
	data := map[string]interface{}{"Eq": eq, "Sort": sort, "OK": true}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditEquipmentHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "equipment_edit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditEquipmentHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var eq Equipment

	eq.Id, _ = strconv.Atoi(r.FormValue("id"))
	eq.Name = r.FormValue("name")
	eq.Type = r.FormValue("type")
	eq.Manufacture = r.FormValue("manufacture")
	eq.Model = r.FormValue("model")
	eq.UniqueNumber = r.FormValue("number")

	fmt.Println(eq)

	err := h.service.UpdateEquipment(context.TODO(), &eq)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/equipment", http.StatusSeeOther)
}
