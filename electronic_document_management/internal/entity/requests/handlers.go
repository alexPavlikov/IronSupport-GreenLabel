package requests

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.ServeFiles("/data/*filepath", http.Dir("data"))

	router.HandlerFunc(http.MethodGet, "/edm/request", h.RequestsHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edit", h.EditRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edits", h.EditPostRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/add", h.InsertRequestHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) RequestsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
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

	data := map[string]interface{}{"Requests": reqs, "RID": RID}
	header := map[string]string{"Title": "ЭДО - Заявки"}
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
	tmpl, err := template.ParseGlob("./internal/html/*.html")
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

	title := map[string]string{"Title": "ЭДО - Изменение заявки"}
	data := map[string]interface{}{"Req": req}

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
