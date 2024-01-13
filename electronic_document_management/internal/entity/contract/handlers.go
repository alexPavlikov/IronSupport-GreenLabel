package contract

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

	// router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/contract", h.ContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/add", h.AddContractHandler)
	router.HandlerFunc(http.MethodPost, "/edm/contract/sorted", h.SortContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/sorted", h.SortContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/edit", h.EditContractHandler)
	router.HandlerFunc(http.MethodGet, "/edm/contract/edits", h.EditPostContractHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ContractHandler(w http.ResponseWriter, r *http.Request) {
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

	title := map[string]string{"Title": "ЭДО - Контракты", "Page": "Contract"}
	data := map[string]interface{}{"Contracts": contracts, "Clients": cls}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "contract", data)
	if err != nil {
		http.NotFound(w, r)
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

		fmt.Println(contracts)

		data := map[string]interface{}{"Contracts": contracts}
		header := map[string]string{"Title": "ЭДО - Контракты", "Page": "Contract"}
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

	title := map[string]string{"Title": "ЭДО - Редактирование контракта", "Page": "Contract"}
	data := map[string]interface{}{"Contract": contract}

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

	fmt.Println("??????????????????????????????????????????????????????????????", contract)

	err := h.service.UpdateContract(context.TODO(), &contract)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/contract", http.StatusSeeOther)
}
