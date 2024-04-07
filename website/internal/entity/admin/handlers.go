package admin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/website/menu", h.WebsiteMenuHandler)
	router.HandlerFunc(http.MethodGet, "/administrator", h.AdminHandler)
	router.HandlerFunc(http.MethodGet, "/administrator/news/edit", h.AdminNewsEdit)
	router.HandlerFunc(http.MethodPost, "/administrator/news/edit", h.AdminNewsEdit)

	router.HandlerFunc(http.MethodGet, "/administrator/news/add", h.AdminNewsAddHandler)
	router.HandlerFunc(http.MethodPost, "/administrator/news/add", h.AdminNewsAddHandler)

	router.HandlerFunc(http.MethodGet, "/administrator/product/add", h.AdminProductAddHandler) //!!

	router.HandlerFunc(http.MethodGet, "/administrator/trust/edit", h.AdminTrustEditHandler)
	router.HandlerFunc(http.MethodPost, "/administrator/trust/edit", h.AdminTrustEditHandler)

	router.HandlerFunc(http.MethodGet, "/administrator/trust/add", h.AdminTrustAddHandler)
	router.HandlerFunc(http.MethodPost, "/administrator/trust/add", h.AdminTrustAddHandler)

	router.HandlerFunc(http.MethodGet, "/administrator/vacancy/edit", h.AdminVacancyEditHandler)
	router.HandlerFunc(http.MethodPost, "/administrator/vacancy/edit", h.AdminVacancyEditHandler)

	router.HandlerFunc(http.MethodGet, "/administrator/vacancy/add", h.AdminVacancyAddHandler)
	router.HandlerFunc(http.MethodPost, "/administrator/vacancy/add", h.AdminVacancyAddHandler)

	router.HandlerFunc(http.MethodPost, "/administrator/email/add", h.AdminEmailAddHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) WebsiteMenuHandler(w http.ResponseWriter, r *http.Request) {
	//if usr.UserAuth.Err {
	tmpl, err := template.ParseGlob("./website/internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open WebsiteMenuHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "menuweb", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open WebsiteMenuHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	//}
}

func (h *handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()

	page := r.FormValue("page")
	if page == "" {
		page = "main"
	}

	news, err := h.service.GetNews(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	product, err := h.service.GetProducts(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	trust, err := h.service.GetTrustCompany(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	vacancy, err := h.service.GetVacancy(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"Body": page, "News": news, "Product": product, "TrustCompany": trust, "Vacancy": vacancy}

	err = tmpl.ExecuteTemplate(w, "panel", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) AdminNewsEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))

		news, err := h.service.GetNewsById(context.TODO(), id)
		if err != nil {
			http.Redirect(w, r, "/administrator", http.StatusSeeOther)
		}

		data := map[string]interface{}{"News": news}

		err = tmpl.ExecuteTemplate(w, "newsedit", data)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		var nw news.News

		r.ParseForm()

		//	nw.Avatar = r.FormValue("avatar")
		//
		file, fileHeader, err := r.FormFile("avatar")
		if err != nil {
			nw.Avatar = r.FormValue("link")
		} else {
			defer file.Close()
			imgPath := fmt.Sprintf("./assets/data/news/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
			dst, err := os.Create(imgPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			nw.Avatar = imgPath[1:]
		}

		nw.Id, _ = strconv.Atoi(r.FormValue("id"))
		nw.Title = r.FormValue("name")
		nw.Text = r.FormValue("text")
		nw.VideoLink = r.FormValue("video")
		d := r.FormValue("deleted")
		if d == "true" {
			nw.Deleted = true
		} else {
			nw.Deleted = false
		}

		h.service.UpdateNews(context.TODO(), nw)

		http.Redirect(w, r, "/administrator?page=news", http.StatusSeeOther)
	}
}

func (h *handler) AdminNewsAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminNewsAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "newsadd", nil)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminNewsAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		var nw news.News

		r.ParseForm()

		file, fileHeader, err := r.FormFile("avatar")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer file.Close()
		imgPath := fmt.Sprintf("./assets/data/news/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		dst, err := os.Create(imgPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//
		nw.Avatar = imgPath[1:]
		slice := utils.ReadCookies(r)
		fmt.Println("#!#2312", slice)
		nw.Author, _ = strconv.Atoi(slice[3])
		nw.CreateDate = time.Now().Format("02.01.2006")
		nw.Title = r.FormValue("name")
		nw.Text = r.FormValue("text")
		nw.VideoLink = r.FormValue("video")
		d := r.FormValue("deleted")
		if d == "true" {
			nw.Deleted = true
		} else {
			nw.Deleted = false
		}

		err = h.service.AddNews(context.TODO(), &nw)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, "/administrator?page=news", http.StatusSeeOther)
		}
	}
}

func (h *handler) AdminProductAddHandler(w http.ResponseWriter, r *http.Request) {
	//add from file with 1c data
}

func (h *handler) AdminTrustEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminTrustEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		r.ParseForm()
		name := r.FormValue("name")

		fmt.Println(name)

		trust, err := h.service.GetTrustCompanyByName(context.TODO(), name)

		fmt.Println(trust, err)

		if err != nil {
			h.logger.Tracef("%s - failed open website AdminTrustEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		data := map[string]interface{}{"Trust": trust}

		err = tmpl.ExecuteTemplate(w, "trustedit", data)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminTrustEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {

		var tc guest.TrustCompany

		r.ParseForm()

		file, fileHeader, err := r.FormFile("logo")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			// fmt.Println(err, "!!!!!!!")
			// return
			tc.Logo = r.FormValue("link")
		} else {
			defer file.Close()
			imgPath := fmt.Sprintf("./assets/data/trust/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)) //
			dst, err := os.Create(imgPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tc.Logo = imgPath
		}

		tc.Name = r.FormValue("name")
		tc.Description = r.FormValue("description")

		err = h.service.UpdateTrustCompany(context.TODO(), tc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/administrator?page=client", http.StatusSeeOther)
	}
}

func (h *handler) AdminTrustAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminTrustAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "trustadd", nil)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminTrustAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		var tc guest.TrustCompany

		r.ParseForm()

		file, fileHeader, err := r.FormFile("logo")
		if err != nil {
			http.NotFound(w, r)
		}
		defer file.Close()
		imgPath := fmt.Sprintf("./assets/data/trust/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)) //
		dst, err := os.Create(imgPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tc.Logo = imgPath
		tc.Name = r.FormValue("name")
		tc.Description = r.FormValue("description")

		err = h.service.AddTrustCompany(context.TODO(), tc)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Redirect(w, r, "/administrator?page=client", http.StatusSeeOther)
		}
	}
}

func (h *handler) AdminVacancyEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		r.ParseForm()
		name := r.FormValue("name")
		vacancy, err := h.service.GetVacancyByName(context.TODO(), name)
		fmt.Println(name, vacancy)
		if err != nil {
			fmt.Println(err)
			h.logger.Tracef("%s - failed open website AdminVacancyEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		data := map[string]interface{}{"Vacancy": vacancy}

		err = tmpl.ExecuteTemplate(w, "vacancyedit", data)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		var v site.Vacancy
		name := r.FormValue("old-name")
		v.Name = r.FormValue("name")
		text := r.FormValue("options")
		b := r.FormValue("active")
		if b == "true" {
			v.Active = false
		} else {
			v.Active = true
		}
		text = text[1 : len(text)-2]
		v.Options = strings.Split(text, " ")

		err := h.service.UpdateVacancy(context.TODO(), v, name)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyEditHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, "/administrator?page=vacancy", http.StatusSeeOther)
		}
	}
}

func (h *handler) AdminVacancyAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseGlob("./website/internal/html/admin/*.html")
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}

		err = tmpl.ExecuteTemplate(w, "vacancyadd", nil)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		var v site.Vacancy
		v.Name = r.FormValue("name")
		text := r.FormValue("options")
		v.Active = true

		v.Options = strings.Split(text, " ")

		err := h.service.AddVacancy(context.TODO(), v)
		if err != nil {
			h.logger.Tracef("%s - failed open website AdminVacancyAddHandler", config.LOG_ERROR)
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, "/administrator?page=vacancy", http.StatusSeeOther)
		}
	}
}

func (h *handler) AdminEmailAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	text := r.FormValue("text")

	sub, err := h.service.GetAllSubscribers(context.TODO())
	if err != nil {
		h.logger.Tracef("%s - failed open website AdminEmailAddHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
	var wg sync.WaitGroup

	for _, v := range sub {
		wg.Add(1)
		go utils.SendMessage(v, text, title)
		defer wg.Done()
	}
	wg.Wait()

	http.Redirect(w, r, "/administrator?page=email", http.StatusSeeOther)
}

// ip, _, err := net.SplitHostPort(r.RemoteAddr)
// if err != nil {
// 	fmt.Fprintf(w, "ip: %q is not IP:port", r.RemoteAddr)
// 	return
// }

// fmt.Println(r.RemoteAddr, ip)

// if ip == `IP который вам нуежен` {
// 	http.Redirect(w, r, `URL 1`, http.StatusFound)
// } else {
// 	http.Redirect(w, r, `URL 2`, http.StatusFound)
// }
