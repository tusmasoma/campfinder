//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tusmasoma/campfinder/db"
)

type SpotCreateRequest struct {
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

type SpotHandler interface {
	HandleSpotCreate(w http.ResponseWriter, r *http.Request)
}

type spotHandler struct {
	sr db.SpotRepository
}

func NewSpotHandler(sr db.SpotRepository) SpotHandler {
	return &spotHandler{
		sr: sr,
	}
}

func (sh *spotHandler) HandleSpotCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody SpotCreateRequest
	if ok := isValidateSpotCreateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := sh.SpotCreate(ctx, requestBody); err != nil {
		http.Error(w, "Internal server error while creating spot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateSpotCreateRequest(body io.ReadCloser, requestBody *SpotCreateRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.Category == "" || requestBody.Name == "" || requestBody.Address == "" || requestBody.Lat == 0 || requestBody.Lng == 0 {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (sh *spotHandler) SpotCreate(ctx context.Context, requestBody SpotCreateRequest) error {
	exists, err := sh.sr.CheckIfSpotExists(ctx, requestBody.Lat, requestBody.Lng)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return err
	}
	if exists {
		log.Printf("User with this name already exists - status: %d", http.StatusConflict)
		return fmt.Errorf("user with this name already exists")
	}

	var spot = db.Spot{
		Category:    requestBody.Category,
		Name:        requestBody.Name,
		Address:     requestBody.Address,
		Lat:         requestBody.Lat,
		Lng:         requestBody.Lng,
		Period:      requestBody.Period,
		Phone:       requestBody.Phone,
		Price:       requestBody.Price,
		Description: requestBody.Description,
		IconPath:    requestBody.IconPath,
	}

	if err = sh.sr.Create(ctx, spot); err != nil {
		log.Printf("Failed to create spot: %v", err)
		return err
	}
	return nil
}
