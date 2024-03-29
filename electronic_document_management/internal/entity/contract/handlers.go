package contract

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

	// router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/contract", h.ContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/add", h.AddContractHandler)
	router.HandlerFunc(http.MethodPost, "/edm/contract/sorted", h.SortContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/sorted", h.SortContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/edit", h.EditContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/edits", h.EditPostContractHandler)

	router.HandlerFunc(http.MethodGet, "/edm/contract/find", h.ContractFindHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var Events []string

func (h *handler) ContractHandler(w http.ResponseWriter, r *http.Request) {

	if !user.UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		contracts, err := h.service.GetContracts(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		cls, err := h.service.GetAllClients(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		Events, err = utils.ReadEventFile()
		if err != nil {
			fmt.Println(err)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Контракты", "Page": "Contract", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"Contracts": contracts, "Clients": cls, "OK": false, "Auth": arr[2]}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "contract", data)
		if err != nil {
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}
}

var contracts []Contract

func (h *handler) SortContractHandler(w http.ResponseWriter, r *http.Request) {
	var ct Contract

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		ct.Client.Name = r.FormValue("client")
		ct.DataStart = r.FormValue("start")
		ct.DataEnd = r.FormValue("end")
		if r.FormValue("status") == "true" {
			ct.Status = true
		} else {
			ct.Status = false
		}

		fmt.Println(ct)

		contracts, err = h.service.GetContractsBySort(context.TODO(), ct)
		if err != nil {
			http.NotFound(w, r)
		}

		http.Redirect(w, r, "/edm/contract/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SortContractHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		arr := utils.ReadCookies(r)

		data := map[string]interface{}{"Contracts": contracts, "OK": true, "Auth": arr[2]}
		header := map[string]interface{}{"Title": "ЭДО - Контракты", "Page": "Contract", "Events": Events, "Auth": arr[2]}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open SortContractHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "contract", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SortContractHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		// err = tmpl.ExecuteTemplate(w, "dialog", dialog)
		// if err != nil {
		// 	h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		// 	//http.NotFound(w, r)
		// }
	}
}

func (h *handler) AddContractHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var contract Contract

	contract.Name = r.FormValue("name")
	contract.Client.Id, _ = strconv.Atoi(r.FormValue("client"))
	contract.DataStart = r.FormValue("startdate")
	contract.DataEnd = r.FormValue("enddate")
	contract.Amount, _ = strconv.Atoi(r.FormValue("price"))
	// contract.File = r.FormValue("file")

	fmt.Println(contract)

	err := h.service.AddContract(context.TODO(), &contract)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/contract", http.StatusSeeOther)
}

func (h *handler) EditContractHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditContractHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	fmt.Println(idx)

	contract, err := h.service.GetContract(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open EditContractHandler GetContract", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	contract.ClientsAll, err = h.service.GetAllClients(context.TODO())
	if err != nil {
		h.logger.Tracef("%s - failed open EditContractHandler GetAllClients", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Редактирование контракта", "Page": "Contract", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Contract": contract, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditContractHandler header", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "contractedit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditContractHandler contractedit", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostContractHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var contract Contract

	contract.Id, _ = strconv.Atoi(r.FormValue("id"))
	contract.Name = r.FormValue("name")
	contract.Client.Id, _ = strconv.Atoi(r.FormValue("client"))
	contract.DataStart = r.FormValue("startdate")
	contract.DataEnd = r.FormValue("enddate")
	contract.Amount, _ = strconv.Atoi(r.FormValue("price"))
	contract.File = r.FormValue("file")
	contract.Status = true

	err := h.service.UpdateContract(context.TODO(), &contract)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/contract", http.StatusSeeOther)
}

func (h *handler) ContractFindHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	text := r.FormValue("text")

	ct, err := h.service.FindContract(context.TODO(), text)
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Поиск", "Page": "Contract", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"Text": text, "Cat": "Contract", "Contracts": ct, "Auth": arr}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "find", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
