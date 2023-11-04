package post

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	router.HandlerFunc(http.MethodGet, "/abs", h.AbsHandler)
	router.HandlerFunc(http.MethodPost, "/abs/move", h.AbsMoveHandler)
	router.HandlerFunc(http.MethodGet, "/abs/delete", h.AbsDeleteHandler)
	router.HandlerFunc(http.MethodGet, "/abs/add", h.AbsAddHandler)

}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) AbsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./abs/internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	posts, err := h.service.GetPosts(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Println(posts)

	data := map[string]interface{}{"Posts": posts}

	err = tmpl.ExecuteTemplate(w, "abs", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) AbsMoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	test := r.FormValue("Ok")
	id := r.FormValue("Id")

	arr := strings.Split(test, " ")
	x := arr[0]
	y := "-" + arr[1]
	x = strings.TrimSuffix(x, "px")
	y = strings.TrimSuffix(y, "px")

	idx, _ := strconv.Atoi(id)
	nx, _ := strconv.Atoi(x)
	ny, _ := strconv.Atoi(y)

	err := h.service.UpdateCordPost(context.TODO(), nx, ny, idx)
	if err != nil {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, "/abs/", http.StatusSeeOther)
	}
}

func (h *handler) AbsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	idx, _ := strconv.Atoi(id)
	err := h.service.DeletePost(context.TODO(), idx)
	if err != nil {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, "/abs/", http.StatusSeeOther)
	}
}

func (h *handler) AbsAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var post Post

	post.Deadline = r.FormValue("date")
	post.Color = fmt.Sprintf(`style="background-color: %s;"`, r.FormValue("color"))
	post.Text = r.FormValue("text")
	post.User, _ = strconv.Atoi(r.FormValue("user"))
	post.PosX = 0
	post.PosY = 0

	err := h.service.AddPost(context.TODO(), &post)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, "/abs/", http.StatusSeeOther)
	}
}
