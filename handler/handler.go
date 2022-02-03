package handler

import (
	"log"
	"net/http"
)

type Handlers struct {
	logger *log.Logger
}

func New(logger *log.Logger) *Handlers {
	return &Handlers{logger: logger}
}

func (h *Handlers) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.audit(h.Root))
}
