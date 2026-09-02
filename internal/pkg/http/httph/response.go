package httph

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/KDarenskii/catalog-service/internal/app/entity"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func SendEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func sendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorApply(r, err)

	var hc httpCoder
	if errors.As(err, &hc) {
		status := hc.HTTPStatus()
		ErrorApplyStatusCode(r, status)
		sendError(w, status, hc)
		return
	}

	ErrorApplyStatusCode(r, http.StatusInternalServerError)
	sendError(w, http.StatusInternalServerError, err)
}

func ParseParam[T any](r *http.Request, name string, parse func(string) (T, error)) (T, error) {
	raw := mux.Vars(r)[name]
	value, err := parse(raw)
	if err != nil {
		var zero T
		return zero, entity.ErrIncorrectParameters
	}
	return value, nil
}

func ParseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	return ParseParam(r, name, uuid.FromString)
}
