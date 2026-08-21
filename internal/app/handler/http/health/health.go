package rhealth

import (
	"log"
	"net/http"

	rhandler "github.com/KDarenskii/catalog-service/internal/app/handler/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("failed to write health check response: %v", err)
	}
}
