package product

import (
	"context"

	"net/http"
	"strconv"
	"text/template"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/products", h.ProductsHandler)
	router.HandlerFunc(http.MethodGet, "/products/sort", h.ProductsSortHandler)
	router.HandlerFunc(http.MethodGet, "/products/product", h.ProductCardHandler)
	router.HandlerFunc(http.MethodGet, "/products/backet/add", h.ProductAddToBacketHandler)
	//router.HandlerFunc(http.MethodPost, "/find", h.FindHandler)
	router.HandlerFunc(http.MethodGet, "/find", h.FindHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

var ProductCat []ProductCategory
var ProductDiscount []DiscountProduct

func (h *handler) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	product, err := h.service.GetProducts(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	ProductCat, err = h.service.GetProductCategory(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	ProductDiscount, err = h.service.GetProductDiscound(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"Product": product, "Category": ProductCat, "Discount": ProductDiscount}
	title := map[string]interface{}{"Guest": guest.Guest, "Title": "Продукция"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "products", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) ProductsSortHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsSortHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	ProductCat, err = h.service.GetProductCategory(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	ProductDiscount, err = h.service.GetProductDiscound(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	r.ParseForm()
	cat := r.FormValue("category")
	price := r.FormValue("price")
	active := r.FormValue("active")
	discount, _ := strconv.Atoi(r.FormValue("discount"))

	product, err := h.service.GetSortedProduct(context.TODO(), cat, price, active, discount)
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsSortHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"Product": product, "Category": ProductCat, "Discount": ProductDiscount}
	title := map[string]interface{}{"Guest": guest.Guest, "Title": "Продукция"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "products", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsSortHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func (h *handler) ProductCardHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//id, _ := strconv.Atoi(r.FormValue("id"))

	// product, err := h.service.GetProductById(context.TODO(), id)
	// if err != nil {
	// 	http.NotFound(w, r)
	// }

	// fmt.Println(product)
}

func (h *handler) ProductAddToBacketHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// guest, _ := strconv.Atoi(r.FormValue("guest"))
	// product, _ := strconv.Atoi(r.FormValue("product"))
}

func (h *handler) FindHandler(w http.ResponseWriter, r *http.Request) {
	var pr []Product
	var err error

	tmpl, err := template.ParseGlob("./website/internal/html/website/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open website FindProductHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	r.ParseForm()

	find := r.FormValue("find")
	if find != "" {
		pr, err = h.service.FindProducts(context.TODO(), find)
		if err != nil {
			http.Redirect(w, r, "/products", http.StatusSeeOther)
		}
	}

	ProductCat, err = h.service.GetProductCategory(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	ProductDiscount, err = h.service.GetProductDiscound(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	data := map[string]interface{}{"Product": pr, "Category": ProductCat, "Discount": ProductDiscount}
	title := map[string]interface{}{"Guest": guest.Guest, "Title": "Продукция"}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		h.logger.Tracef("%s - failed open website IndexHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "products", data)
	if err != nil {
		h.logger.Tracef("%s - failed open website ProductsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

}
