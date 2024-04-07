//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

type SpotHandler interface {
	CreateSpot(w http.ResponseWriter, r *http.Request)
	ListSpots(w http.ResponseWriter, r *http.Request)
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

type ListSpotsResponse struct {
	Spots []model.Spot `json:"spots"`
}

func (sh *spotHandler) CreateSpot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody CreateSpotRequest
	if ok := isValidateCreateSpotRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := sh.suc.CreateSpot(
		ctx,
		requestBody.Category,
		requestBody.Name,
		requestBody.Address,
		requestBody.Lat,
		requestBody.Lng,
		requestBody.Period,
		requestBody.Phone,
		requestBody.Price,
		requestBody.Description,
		requestBody.IconPath,
	)
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

func (sh *spotHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories := r.URL.Query()["category"]
	// spotID := r.URL.Query().Get("spot_id")

	allSpots := sh.suc.ListSpots(ctx, categories)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ListSpotsResponse{Spots: allSpots}); err != nil {
		http.Error(w, "Failed to encode spots to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
