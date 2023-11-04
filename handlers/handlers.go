package handlers

import "github.com/julienschmidt/httprouter"

type Handlers interface {
	Register(router *httprouter.Router)
}
