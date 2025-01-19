package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var request model.ReadTODORequest

		params := r.URL.Query()

		if prevId := params.Get("prev_id"); prevId != "" {
			prevId, err := strconv.ParseInt(prevId, 10, 64)
			if err != nil {
				log.Println(err)
				http.Error(w, "400 BadRequest", http.StatusBadRequest)
				return
			}
			request.PrevID = prevId
		}

		if size := params.Get("size"); size != "" {
			size, err := strconv.ParseInt(size, 10, 64)
			if err != nil {
				log.Println(err)
				http.Error(w, "400 BadRequest", http.StatusBadRequest)
				return
			}
			request.Size = size
		}

		response, err := h.Read(r.Context(), &request)
		if err != nil {
			log.Println(err)
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var request model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println(err)
			http.Error(w, "400 Badrequest", http.StatusBadRequest)
			return
		}

		if request.Subject == "" {
			http.Error(w, "400 BadRequest", http.StatusBadRequest)
			return
		}

		response, err := h.Create(r.Context(), &request)
		if err != nil {
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		var request model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println(err)
			http.Error(w, "400 BadRequest", http.StatusBadRequest)
			return
		}

		if request.ID == 0 || request.Subject == "" {
			http.Error(w, "400 BadRequest", http.StatusBadRequest)
			return
		}

		response, err := h.Update(r.Context(), &request)
		if err != nil {
			log.Println(err)

			if _, ok := err.(*model.ErrNotFound); ok {
				http.Error(w, "404 NotFound", http.StatusNotFound)
			} else {
				http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		var request model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println(err)
			http.Error(w, "400 BadRequest", http.StatusBadRequest)
			return
		}

		if len(request.IDs) == 0 {
			http.Error(w, "400 BadRequest", http.StatusBadRequest)
			return
		}

		response, err := h.Delete(r.Context(), &request)
		if err != nil {
			log.Println(err)

			if _, ok := err.(*model.ErrNotFound); ok {
				http.Error(w, "404 NotFound", http.StatusNotFound)
			} else {
				http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)

	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)

	return &model.ReadTODOResponse{TODOs: todos}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)

	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)

	return &model.DeleteTODOResponse{}, err
}
