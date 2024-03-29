package stat

import (
	"fmt"
	"net/http"
	"text/template"

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

// Register implements handlers.Handlers.
func (h *handler) Register(router *httprouter.Router) {

	router.HandlerFunc(http.MethodGet, "/edm/statistics", h.StatHandler)

	// router.HandlerFunc(http.MethodGet, "/edm/user", h.UserHandler)
	// router.HandlerFunc(http.MethodPost, "/edm/user/sorted", h.SortUserHandler)
	// router.HandlerFunc(http.MethodGet, "/edm/user/sorted", h.SortUserHandler)
	// router.HandlerFunc(http.MethodGet, "/edm/user/add", h.AddUserHandler)
	// router.HandlerFunc(http.MethodGet, "/edm/user/edit", h.EditUserHandler)
	// router.HandlerFunc(http.MethodGet, "/edm/user/edits", h.EditPostUserHandler)
	// router.HandlerFunc(http.MethodGet, "/edm/user/account", h.AccountHandler)

	// router.HandlerFunc(http.MethodGet, "/edm/user/role/add", h.AddRoleUserHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var Events []string

// func (h *handler) UserHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	users, err := h.service.GetUsers(context.TODO())
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	role, err := h.service.GetRole(context.TODO())
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	fmt.Println(users)

// 	title := map[string]string{"Title": "ЭДО - Пользователи", "Page": "User"}
// 	data := map[string]interface{}{"User": users, "Role": role}

// 	err = tmpl.ExecuteTemplate(w, "header", title)
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	err = tmpl.ExecuteTemplate(w, "user", data)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.NotFound(w, r)
// 	}
// }

// var users []User

// func (h *handler) SortUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	r.ParseForm()

// 	user.FullName = r.FormValue("fio")
// 	user.Email = r.FormValue("email")
// 	user.Phone = r.FormValue("phone")
// 	user.Role = r.FormValue("role")

// 	var err error

// 	if r.Method == "POST" {
// 		r.ParseForm()

// 		fmt.Println(user)

// 		if user.FullName != "" {
// 			user.FullName = "%" + user.FullName + "%"
// 		}
// 		if user.Email != "" {
// 			user.Email = "%" + user.Email + "%"
// 		}
// 		if user.Phone != "" {
// 			user.Phone = "%" + user.Phone + "%"
// 		}

// 		users, err = h.service.GetUserBySort(context.TODO(), &user)
// 		if err != nil {
// 			http.NotFound(w, r)
// 		}
// 		fmt.Println(users)

// 		http.Redirect(w, r, "/edm/user/sorted", http.StatusSeeOther)

// 	} else if r.Method == "GET" {

// 		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
// 		if err != nil {
// 			fmt.Println(err)
// 			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
// 			w.WriteHeader(http.StatusNotFound)
// 		}

// 		fmt.Println(users)

// 		role, err := h.service.GetRole(context.TODO())
// 		if err != nil {
// 			http.NotFound(w, r)
// 		}

// 		data := map[string]interface{}{"User": users, "Role": role}
// 		header := map[string]string{"Title": "ЭДО - Сотрудники", "Page": "User"}
// 		// dialog := map[string]interface{}{"ReqInsertData": RID}

// 		err = tmpl.ExecuteTemplate(w, "header", header)
// 		if err != nil {
// 			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
// 			//http.NotFound(w, r)
// 		}

// 		err = tmpl.ExecuteTemplate(w, "user", data)
// 		if err != nil {
// 			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
// 			//http.NotFound(w, r)
// 		}
// 	}
// }

// func (h *handler) AddUserHandler(w http.ResponseWriter, r *http.Request) {

// 	r.ParseForm()

// 	var user User

// 	user.FullName = r.FormValue("fio")
// 	user.Email = r.FormValue("email")
// 	user.Phone = r.FormValue("phone")
// 	user.Role = r.FormValue("role")

// 	fmt.Println(user)

// 	err := h.service.AddUser(context.TODO(), &user)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.NotFound(w, r)
// 	}

// 	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
// }

// func (h *handler) AddRoleUserHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	name := r.FormValue("name")
// 	err := h.service.AddRole(context.TODO(), name)
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
// }

// func (h *handler) EditUserHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
// 	if err != nil {
// 		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
// 		http.NotFound(w, r)
// 	}

// 	r.ParseForm()
// 	id := r.FormValue("id")
// 	idx, _ := strconv.Atoi(id)
// 	fmt.Println(idx)

// 	us, err := h.service.GetUser(context.TODO(), idx)
// 	if err != nil {
// 		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
// 		http.NotFound(w, r)
// 	}

// 	role, err := h.service.GetRole(context.TODO())
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	title := map[string]string{"Title": "ЭДО - Редактирование сотрудника", "Page": "User"}
// 	data := map[string]interface{}{"User": us, "Role": role}

// 	err = tmpl.ExecuteTemplate(w, "header", title)
// 	if err != nil {
// 		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
// 		http.NotFound(w, r)
// 	}
// 	err = tmpl.ExecuteTemplate(w, "user_edit", data)
// 	if err != nil {
// 		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
// 		http.NotFound(w, r)
// 	}
// }

// func (h *handler) EditPostUserHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	var us User

// 	us.Id, _ = strconv.Atoi(r.FormValue("id"))
// 	us.FullName = r.FormValue("fio")
// 	us.Email = r.FormValue("email")
// 	us.Phone = r.FormValue("phone")
// 	us.Image = r.FormValue("avatar")
// 	us.Role = r.FormValue("role")

// 	fmt.Println(us)

// 	err := h.service.UpdateUser(context.TODO(), &us)
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
// }

// func (h *handler) AccountHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}

// 	title := map[string]string{"Title": "ЭДО - Пользователи", "Page": "User"}
// 	data := map[string]interface{}{"User": users}

// 	err = tmpl.ExecuteTemplate(w, "header", title)
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}
// 	err = tmpl.ExecuteTemplate(w, "account", data)
// 	if err != nil {
// 		http.NotFound(w, r)
// 	}
// }

func (h *handler) StatHandler(w http.ResponseWriter, r *http.Request) {
	if !user.UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		Events, err = utils.ReadEventFile()
		if err != nil {
			fmt.Println(err)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Статистика", "Page": "Statistics", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}
		err = tmpl.ExecuteTemplate(w, "statistics", data)
		if err != nil {
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}
}
