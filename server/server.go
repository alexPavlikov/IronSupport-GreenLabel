package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"text/template"
	"time"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/handlers"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

var ClientPostgreSQL dbClient.Client

func Start(r *httprouter.Router, cfg config.Config) {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Start application")

	var listener net.Listener
	var listenerErr error

	logger.Info(config.LOG_INFO, "Listen TCP")
	listener, listenerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Infof("%s Server is listen on port: %s:%s", config.LOG_INFO, cfg.Listen.BindIP, cfg.Listen.Port)
	if listenerErr != nil {
		logger.Fatal(config.LOG_ERROR, listenerErr.Error())
	}

	server := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	err := server.Serve(listener)
	if err != nil {
		logger.Fatal(config.LOG_ERROR, err.Error())
	}

}

type handler struct {
	logger *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	router.ServeFiles("/assets/*filepath", http.Dir("./assets/")) //
	router.ServeFiles("/data/*filepath", http.Dir("./data/"))

	router.HandlerFunc(http.MethodGet, "/IronSupport", h.ISHandler)
}

func NewHandler(logger *logging.Logger) handlers.Handlers {
	return &handler{
		logger: logger,
	}
}

func (h *handler) ISHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open ISHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "isgl", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		http.NotFound(w, r)
	}
}

func init() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()

	ClientPostgreSQL, err := dbClient.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("failed to get new client postgresql, due to err: %v", err)
	}
	fmt.Println(ClientPostgreSQL)
}
