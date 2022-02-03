package handler

import (
	"net/http"
)

const msg = `{"msg":"hello world"}`

func (h *Handlers) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
