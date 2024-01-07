package requests

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

	// router.ServeFiles("/assets/*filepath", http.Dir("./assets/"))
	// router.ServeFiles("/data/*filepath", http.Dir("./data/"))

	router.HandlerFunc(http.MethodGet, "/edm/request", h.RequestsHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edit", h.EditRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edits", h.EditPostRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/add", h.InsertRequestHandler)
	router.HandlerFunc(http.MethodPost, "/edm/request/sorted", h.SorterRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/sorted", h.SorterRequestHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) RequestsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		fmt.Println(err)
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}

	reqs, err := h.service.GetRequests(context.TODO())
	if err != nil {
		h.logger.Errorf("%s - failed load RequestsHandler due to err: %s", config.LOG_ERROR, err)
		http.NotFound(w, r)
	}

	fmt.Println(reqs[0].ClientObject.Object.Id, reqs[0].ClientObject.Object.Name)

	data := map[string]interface{}{"Requests": reqs, "RID": RID}
	header := map[string]string{"Title": "ЭДО - Заявки", "Page": "Request"}
	dialog := map[string]interface{}{"ReqInsertData": RID}

	err = tmpl.ExecuteTemplate(w, "header", header)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		//http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "request", data)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		//http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "dialog", dialog)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		//http.NotFound(w, r)
	}
}

func (h *handler) EditRequestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditRequestHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)

	req, err := h.service.GetRequest(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open EditRequestHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	fmt.Println(req)

	title := map[string]string{"Title": "ЭДО - Изменение заявки", "Page": "Request"}
	data := map[string]interface{}{"Req": req, "List": RID}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditRequestHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "regedit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditRequestHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostRequestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var req Request
	id := r.FormValue("id")
	req.Id, _ = strconv.Atoi(id)
	req.Title = r.FormValue("type")
	req.Status.Name = r.FormValue("status")

	//

	req.Priority = r.FormValue("priority")
	req.Description = r.FormValue("description")
	clientId := r.FormValue("client")
	req.Client.Id, _ = strconv.Atoi(clientId)
	contractId := r.FormValue("contract")
	req.Contract.Id, _ = strconv.Atoi(contractId)

	//

	objectId := r.FormValue("object")
	req.ClientObject.Object.Id, _ = strconv.Atoi(objectId)
	fmt.Println(objectId, req.ClientObject.Object.Id)
	//

	equipmentId := r.FormValue("equipment")
	req.Equipment.Id, _ = strconv.Atoi(equipmentId)

	//

	workerId := r.FormValue("worker")
	req.Worker.Id, _ = strconv.Atoi(workerId)

	//

	req.EndDate = r.FormValue("date")

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", req.Priority)

	err := h.service.UpdateRequest(context.TODO(), &req)
	if err != nil {
		http.NotFound(w, r)
		fmt.Println(err)
	}

	// fmt.Println(req)

	// nReq, err := h.service.GetRequest(context.TODO(), req.Id)
	// if err != nil {
	// 	http.NotFound(w, r)
	// 	fmt.Println(err)
	// }

	// fmt.Println(nReq)

	http.Redirect(w, r, "/edm/request", http.StatusSeeOther)
}

func (h *handler) InsertRequestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var req Request
	id := r.FormValue("id")
	req.Id, _ = strconv.Atoi(id)
	req.Title = r.FormValue("type")
	req.Status.Name = r.FormValue("status")

	req.Priority = r.FormValue("priority")

	req.Description = r.FormValue("description")

	clientId := r.FormValue("client")
	req.Client.Id, _ = strconv.Atoi(clientId)

	contractId := r.FormValue("contract")
	req.Contract.Id, _ = strconv.Atoi(contractId)

	objectId := r.FormValue("object")
	req.ClientObject.Object.Id, _ = strconv.Atoi(objectId)

	equipmentId := r.FormValue("equipment")
	req.Equipment.Id, _ = strconv.Atoi(equipmentId)

	workerId := r.FormValue("worker")
	req.Worker.Id, _ = strconv.Atoi(workerId)

	req.EndDate = r.FormValue("date")

	err := h.service.AddRequest(context.TODO(), &req)
	if err != nil {
		http.NotFound(w, r)
	}
}

var rs []Request

func (h *handler) SorterRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req Request

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		req.Client.Id, _ = strconv.Atoi(r.FormValue("client-sort"))
		req.Worker.Id, _ = strconv.Atoi(r.FormValue("worker-sort"))
		req.ClientObject.Id, _ = strconv.Atoi(r.FormValue("object-sort"))
		req.Equipment.Id, _ = strconv.Atoi(r.FormValue("equipment-sort"))
		req.Status.Name = r.FormValue("status-sort")

		rs, err = h.service.GetRequestsBySort(context.TODO(), req)
		if err != nil {
			http.NotFound(w, r)
		}

		http.Redirect(w, r, "/edm/request/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		data := map[string]interface{}{"Requests": rs, "RID": RID}
		header := map[string]string{"Title": "ЭДО - Заявки", "Page": "Request"}
		dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "request", data)
		if err != nil {
			h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "dialog", dialog)
		if err != nil {
			h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}
	}

}
