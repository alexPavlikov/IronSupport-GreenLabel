package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

var UserAuth Auth

// Register implements handlers.Handlers.
func (h *handler) Register(router *httprouter.Router) {

	UserAuth.Err = true

	router.HandlerFunc(http.MethodGet, "/edm/user", h.UserHandler)
	router.HandlerFunc(http.MethodPost, "/edm/user/sorted", h.SortUserHandler)
	router.HandlerFunc(http.MethodGet, "/edm/user/sorted", h.SortUserHandler)
	router.HandlerFunc(http.MethodGet, "/edm/user/add", h.AddUserHandler)
	router.HandlerFunc(http.MethodGet, "/edm/user/edit", h.EditUserHandler)
	router.HandlerFunc(http.MethodGet, "/edm/user/edits", h.EditPostUserHandler)
	router.HandlerFunc(http.MethodGet, "/edm/user/account", h.AccountHandler)

	router.HandlerFunc(http.MethodGet, "/edm/user/role/add", h.AddRoleUserHandler)

	router.HandlerFunc(http.MethodGet, "/user/auth", h.UserAuthHandler)
	router.HandlerFunc(http.MethodGet, "/user/auth/authconfirm", h.AuthUserConfirm)
	router.HandlerFunc(http.MethodGet, "/user/auth/reg-confirm", h.RegUserConfirm)
	router.HandlerFunc(http.MethodGet, "/user/exit", h.UserExitHandler)

	router.HandlerFunc(http.MethodGet, "/edm/user/find", h.UserFindHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var Events []string

func (h *handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	if !UserAuth.Err {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			http.NotFound(w, r)
		}

		users, err := h.service.GetUsers(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		role, err := h.service.GetRole(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		Events, err = utils.ReadEventFile()
		if err != nil {
			fmt.Println(err)
		}

		arr := utils.ReadCookies(r)

		title := map[string]interface{}{"Title": "ЭДО - Пользователи", "Page": "User", "Events": Events, "Auth": arr[2]}
		data := map[string]interface{}{"User": users, "Role": role, "OK": false, "Auth": arr[2]}

		err = tmpl.ExecuteTemplate(w, "header", title)
		if err != nil {
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "user", data)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
	} else {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}
}

var users []User

func (h *handler) SortUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	r.ParseForm()

	user.FullName = r.FormValue("fio")
	user.Email = r.FormValue("email")
	user.Phone = r.FormValue("phone")
	user.Role = r.FormValue("role")

	var err error

	if r.Method == "POST" {
		r.ParseForm()

		fmt.Println(user)

		if user.FullName != "" {
			user.FullName = "%" + user.FullName + "%"
		}
		if user.Email != "" {
			user.Email = "%" + user.Email + "%"
		}
		if user.Phone != "" {
			user.Phone = "%" + user.Phone + "%"
		}

		users, err = h.service.GetUserBySort(context.TODO(), &user)
		if err != nil {
			http.NotFound(w, r)
		}
		fmt.Println(users)

		http.Redirect(w, r, "/edm/user/sorted", http.StatusSeeOther)

	} else if r.Method == "GET" {

		tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
			w.WriteHeader(http.StatusNotFound)
		}

		fmt.Println(users)

		role, err := h.service.GetRole(context.TODO())
		if err != nil {
			http.NotFound(w, r)
		}

		arr := utils.ReadCookies(r)

		data := map[string]interface{}{"User": users, "Role": role, "OK": true, "Auth": arr[2]}
		header := map[string]interface{}{"Title": "ЭДО - Сотрудники", "Page": "User", "Events": Events, "Auth": arr[2]}
		// dialog := map[string]interface{}{"ReqInsertData": RID}

		err = tmpl.ExecuteTemplate(w, "header", header)
		if err != nil {
			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "user", data)
		if err != nil {
			h.logger.Tracef("%s - failed open SortUserHandler", config.LOG_ERROR)
			//http.NotFound(w, r)
		}
	}
}

func (h *handler) AddUserHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var user User

	user.FullName = r.FormValue("fio")
	user.Email = r.FormValue("email")
	user.Phone = r.FormValue("phone")
	user.Role = r.FormValue("role")

	fmt.Println(user)

	err := h.service.AddUser(context.TODO(), &user)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
}

func (h *handler) AddRoleUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	err := h.service.AddRole(context.TODO(), name)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
}

func (h *handler) EditUserHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	fmt.Println(idx)

	us, err := h.service.GetUser(context.TODO(), idx)
	if err != nil {
		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	role, err := h.service.GetRole(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)

	title := map[string]interface{}{"Title": "ЭДО - Редактирование сотрудника", "Page": "User", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"User": us, "Role": role, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "user_edit", data)
	if err != nil {
		h.logger.Tracef("%s - failed open EditUserHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) EditPostUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var us User

	us.Id, _ = strconv.Atoi(r.FormValue("id"))
	us.FullName = r.FormValue("fio")
	us.Email = r.FormValue("email")
	us.Phone = r.FormValue("phone")
	us.Image = r.FormValue("avatar")
	us.Role = r.FormValue("role")

	fmt.Println(us)

	err := h.service.UpdateUser(context.TODO(), &us)
	if err != nil {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, "/edm/user", http.StatusSeeOther)
}

func (h *handler) AccountHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	arr := utils.ReadCookies(r)
	user, err := h.service.GetAuthUser(context.TODO(), arr[0], arr[1])
	if err != nil {
		http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
	}

	title := map[string]interface{}{"Title": "ЭДО - Пользователи", "Page": "User", "Events": Events, "Auth": arr[2]}
	data := map[string]interface{}{"User": user, "Auth": arr[2]}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "account", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()
	v := r.FormValue("val")

	text := "singup"
	if v != "" {
		text = "reg"
	}

	// title := map[string]interface{}{}
	data := map[string]interface{}{"Content": text, "UserAuth": UserAuth, "Title": "Авторизация"}

	err = tmpl.ExecuteTemplate(w, "auth", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) AuthUserConfirm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var u User
	var err error

	u.Email = r.FormValue("email")
	u.Password = r.FormValue("pass")

	if u.Email != "" && u.Password != "" {

		fmt.Println("read form auth", u)

		UserAuth.Us, err = h.service.GetAuthUser(context.TODO(), u.Email, u.Password)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("GetAuthUser", UserAuth.Us, UserAuth.Err)

		if err != nil {
			fmt.Println(err)
			UserAuth.Err = true
			http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
		} else if UserAuth.Us.Password == u.Password {
			UserAuth.Err = false

			expires := time.Now().AddDate(1, 0, 0)
			cookie := &http.Cookie{
				Name:  "Id",
				Value: UserAuth.Us.Email + " " + UserAuth.Us.Password + " " + UserAuth.Us.Role + " " + fmt.Sprint(UserAuth.Us.Id),
				//MaxAge:  300,
				Expires: expires,
				Path:    "/",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/IronSupport", http.StatusSeeOther)

		} else if UserAuth.Us.Password != u.Password {
			UserAuth.Err = true
			http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
		}
	}

}

func (h *handler) RegUserConfirm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var u User

	u.FullName = r.FormValue("fullname")
	u.Email = r.FormValue("email")
	u.Phone = r.FormValue("phone")
	u.Password = r.FormValue("pass")
	// pass := r.FormValue("xpass")

}

func (h *handler) UserFindHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./electronic_document_management/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()

	text := r.FormValue("text")

	us, err := h.service.FindUser(context.TODO(), text)
	if err != nil {
		http.NotFound(w, r)
	}

	title := map[string]interface{}{"Title": "ЭДО - Поиск", "Page": "User", "Events": Events}
	data := map[string]interface{}{"Text": text, "Cat": "User", "User": us}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}
	err = tmpl.ExecuteTemplate(w, "find", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) UserExitHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "Id",
		Value: "",
		//MaxAge:  300,
		Expires: time.Unix(0, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/user/auth", http.StatusSeeOther)
}
