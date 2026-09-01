package hcategory

import (
	"net/http"

	"github.com/KDarenskii/catalog-service/internal/app/entity"
	rhandler "github.com/KDarenskii/catalog-service/internal/app/handler/http"
	"github.com/KDarenskii/catalog-service/internal/app/service"
	"github.com/KDarenskii/catalog-service/internal/pkg/http/binding"
	"github.com/KDarenskii/catalog-service/internal/pkg/http/httph"
)

type handler struct {
	srv service.Category
}

func NewHandler(srv service.Category) rhandler.Category {
	return &handler{srv: srv}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestCategoryCreate

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	category, err := h.srv.Create(r.Context(), req)
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	resp := entity.ResponseCategoryCreate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}

	httph.SendJSON(w, http.StatusCreated, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	guid, err := httph.ParseUUIDParam(r, "guid")
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	var req entity.RequestCategoryUpdate

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	category, err := h.srv.Update(r.Context(), guid, req)
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	resp := entity.ResponseCategoryUpdate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httph.SendJSON(w, http.StatusOK, resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	guid, err := httph.ParseUUIDParam(r, "guid")
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	err = h.srv.Delete(r.Context(), guid)
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	httph.SendEmpty(w, http.StatusOK)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.srv.List(r.Context())
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	responseItems := make([]entity.ResponseCategoryListItem, 0, len(categories))

	for _, category := range categories {
		responseItems = append(responseItems, entity.ResponseCategoryListItem{
			GUID:      category.GUID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}

	httph.SendJSON(w, http.StatusOK, entity.ResponseCategoryList{Data: responseItems})
}
