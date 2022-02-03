package handler

import (
	"net/http"
	"time"
)

func (h *Handlers) audit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("[INFO] request processed at / in %v\n", time.Since(startTime))
		next(w, r)
	}
}
