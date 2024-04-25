//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

type SpotHandler interface {
	CreateSpot(w http.ResponseWriter, r *http.Request)
	BatchCreateSpots(w http.ResponseWriter, r *http.Request)
	ListSpots(w http.ResponseWriter, r *http.Request)
	GetSpot(w http.ResponseWriter, r *http.Request)
}

type spotHandler struct {
	suc usecase.SpotUseCase
}

func NewSpotHandler(suc usecase.SpotUseCase) SpotHandler {
	return &spotHandler{
		suc: suc,
	}
}

type CreateSpotRequest struct {
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Period      string  `json:"period"`
	Phone       string  `json:"phone"`
	Price       string  `json:"price"`
	Description string  `json:"description"`
	IconPath    string  `json:"iconpath"`
}

type BatchCreateSpotsRequest struct {
	Spots []CreateSpotRequest `json:"spots"`
}

type ListSpotsResponse struct {
	Spots []model.Spot `json:"spots"`
}

type GetSpotResponse struct {
	Spot model.Spot `json:"spot"`
}

func (sh *spotHandler) CreateSpot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody CreateSpotRequest

	defer r.Body.Close()
	if ok := isValidateCreateSpotRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}

	params := convertCreateSpotRequestToParams(requestBody)
	err := sh.suc.CreateSpot(ctx, &params)
	if err != nil {
		http.Error(w, "Internal server error while creating spot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateCreateSpotRequest(body io.ReadCloser, requestBody *CreateSpotRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.Category == "" ||
		requestBody.Name == "" ||
		requestBody.Address == "" ||
		requestBody.Lat == 0 ||
		requestBody.Lng == 0 {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func convertCreateSpotRequestToParams(req CreateSpotRequest) usecase.CreateSpotParams {
	return usecase.CreateSpotParams{
		Category:    req.Category,
		Name:        req.Name,
		Address:     req.Address,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Period:      req.Period,
		Phone:       req.Phone,
		Price:       req.Price,
		Description: req.Description,
		IconPath:    req.IconPath,
	}
}

func (sh *spotHandler) BatchCreateSpots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody BatchCreateSpotsRequest

	defer r.Body.Close()
	if ok := isValidateBatchCreateSpotsRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user batch create request", http.StatusBadRequest)
		return
	}

	params := convertBatchCreateSpotsRequestToParams(requestBody)
	err := sh.suc.BatchCreateSpots(ctx, &params)
	if err != nil {
		http.Error(w, "Internal server error while batch creating spots", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateBatchCreateSpotsRequest(body io.ReadCloser, requestBody *BatchCreateSpotsRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	for _, spot := range requestBody.Spots {
		if spot.Category == "" ||
			spot.Name == "" ||
			spot.Address == "" ||
			spot.Lat == 0 ||
			spot.Lng == 0 {
			log.Printf("Missing required fields")
			return false
		}
	}
	return true
}

func convertBatchCreateSpotsRequestToParams(req BatchCreateSpotsRequest) usecase.BatchCreateSpotParams {
	var params []usecase.CreateSpotParams
	for _, spotReq := range req.Spots {
		spotParam := convertCreateSpotRequestToParams(spotReq)
		params = append(params, spotParam)
	}
	return usecase.BatchCreateSpotParams{
		Spots: params,
	}
}

func (sh *spotHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories := r.URL.Query()["category"]

	allSpots := sh.suc.ListSpots(ctx, categories)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ListSpotsResponse{Spots: allSpots}); err != nil {
		http.Error(w, "Failed to encode spots to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *spotHandler) GetSpot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotID := chi.URLParam(r, "spotID")

	spot := sh.suc.GetSpot(ctx, spotID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(GetSpotResponse{Spot: spot}); err != nil {
		http.Error(w, "Failed to encode spot to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
