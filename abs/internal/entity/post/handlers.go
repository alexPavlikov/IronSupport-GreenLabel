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
	router.HandlerFunc(http.MethodPost, "/abs/edit", h.EditAbsHandler)

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

	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("Id"))
	editPost, err := h.service.GetPost(context.TODO(), id)
	if err != nil {
		fmt.Println(err)
	}

	editArr := []Post{editPost}

	data := map[string]interface{}{"Posts": posts, "EditPost": editArr}

	err = tmpl.ExecuteTemplate(w, "abs", data)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *handler) AbsMoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	test := r.FormValue("Position")
	id := r.FormValue("Id")

	arr := strings.Split(test, " ")
	if len(arr) > 2 {
		x := arr[0]
		y := arr[3]
		x = strings.TrimSuffix(x, "px")
		y = strings.TrimSuffix(y, "px")

		idx, _ := strconv.Atoi(id)
		nx, err := strconv.Atoi(x)
		if err != nil {
			fmt.Println(err)
		}
		ny, err := strconv.Atoi(y)
		if err != nil {
			fmt.Println(err)
		}

		err = h.service.UpdateCordPost(context.TODO(), nx, ny, idx)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, "/abs/", http.StatusSeeOther)
		}
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

func (h *handler) EditAbsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var post Post

	post.Id, _ = strconv.Atoi(r.FormValue("id"))
	post.Text = r.FormValue("text")
	fmt.Println(post)
	err := h.service.repository.UpdateTextAbs(context.TODO(), post.Text, post.Id)
	if err != nil {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, "/abs/", http.StatusSeeOther)
	}
}
