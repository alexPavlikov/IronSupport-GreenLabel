package requests

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

var Events []string

func (h *handler) Register(router *httprouter.Router) {

	// router.ServeFiles("/assets/*filepath", http.Dir("./assets/"))
	// router.ServeFiles("/data/*filepath", http.Dir("./data/"))

	router.HandlerFunc(http.MethodGet, "/edm/request", h.RequestsHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edit", h.EditRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/edits", h.EditPostRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/add", h.InsertRequestHandler)
	router.HandlerFunc(http.MethodPost, "/edm/request/sorted", h.SorterRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/sorted", h.SorterRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/answer", h.AnswerRequestHandler)
	router.HandlerFunc(http.MethodPost, "/edm/request/answer", h.AnswerRequestHandler)
	router.HandlerFunc(http.MethodGet, "/edm/request/find", h.RequestFindHandler)

	router.HandlerFunc(http.MethodGet, "/edm/find", h.FindHandler)

	var err error
	Events, err = utils.ReadEventFile() // возможно сделать во всех папках
	if err != nil {
		fmt.Println(err)
	}
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) RequestsHandler(w http.ResponseWriter, r *http.Request) {
	if !user.UserAuth.Err {
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

		arr := utils.ReadCookies(r)

		data := map[string]interface{}{"Requests": reqs, "RID": RID, "OK": false}
		header := map[string]interface{}{"Title": "ЭДО - Заявки", "Page": "Request", "Events": Events, "Auth": arr[2]}
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
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
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

	arr := utils.ReadCookies(r)

	us, err := h.service.repository.GetRequestWorkerByEmail(context.TODO(), arr[0])
	if err != nil {
		h.logger.Tracef("%s - failed open EditRequestHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	title := map[string]interface{}{"Title": "ЭДО - Изменение заявки", "Page": "Request", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Req": req, "List": RID, "Auth": us}

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

	req.Name = r.FormValue("name")

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

	http.Redirect(w, r, "/edm/request", http.StatusSeeOther)

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

		arr := utils.ReadCookies(r)

		data := map[string]interface{}{"Requests": rs, "RID": RID, "OK": true}
		header := map[string]interface{}{"Title": "ЭДО - Заявки", "Page": "Request", "Events": Events, "Auth": arr[2]}
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

func (h *handler) AnswerRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		r.ParseForm()

		id, _ := strconv.Atoi(r.FormValue("id"))

		req, err := h.service.GetRequest(context.TODO(), id)
		if err != nil {
			http.NotFound(w, r)
		}

		ra, err := h.service.GetAnswerRequest(context.TODO(), id)
		if err != nil {
			http.NotFound(w, r)
		}

		arr := utils.ReadCookies(r)
		user, err := h.service.repository.GetRequestWorkerByEmail(context.TODO(), arr[0])
		if err != nil {
			http.NotFound(w, r)
		}

		title := map[string]interface{}{"Title": "ЭДО - Заявка", "Page": "Request", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"Req": req, "Answer": ra, "User": user}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "answer", data)
		if err != nil {
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		var ra ReqAns

		ra.Text = r.FormValue("text")
		ra.Request.Id, _ = strconv.Atoi(r.FormValue("request"))
		ra.Worker.Id, _ = strconv.Atoi(r.FormValue("user"))

		err := h.service.AddAnswerRequest(context.TODO(), &ra)
		if err != nil {
			http.NotFound(w, r)
		}

		link := "/edm/request/answer?id=" + fmt.Sprint(ra.Request.Id)
		http.Redirect(w, r, link, http.StatusSeeOther)

	} else {
		http.NotFound(w, r)
	}
}

func (h *handler) FindHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	category := r.FormValue("category")
	text := r.FormValue("text")

	switch category {
	case "Request":
		http.Redirect(w, r, "/edm/request/find?text="+text, http.StatusSeeOther)
	case "Client":
		http.Redirect(w, r, "/edm/client/find?text="+text, http.StatusSeeOther)
	case "Contract":
		http.Redirect(w, r, "/edm/contract/find?text="+text, http.StatusSeeOther)
	case "Object":
		http.Redirect(w, r, "/edm/object/find?text="+text, http.StatusSeeOther)
	case "Equipment":
		http.Redirect(w, r, "/edm/equipment/find?text="+text, http.StatusSeeOther)
	case "Worker":
		http.Redirect(w, r, "/edm/user/find?text="+text, http.StatusSeeOther)
	case "Service":
		http.Redirect(w, r, "/edm/service/find?text="+text, http.StatusSeeOther)

	}

}

func (h *handler) RequestFindHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	text := r.FormValue("text")

	req, err := h.service.FindRequest(context.TODO(), text)
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Поиск", "Page": "Request", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Text": text, "Cat": "Request", "Req": req}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "find", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
