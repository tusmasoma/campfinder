//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

const (
	DefaultUUID = "00000000-0000-0000-0000-000000000000"
)

type ImageHandler interface {
	HandleImageGet(w http.ResponseWriter, r *http.Request)
	HandleImageCreate(w http.ResponseWriter, r *http.Request)
	HandleImageDelete(w http.ResponseWriter, r *http.Request)
}

type imageHandler struct {
	iuc usecase.ImageUseCase
	auc usecase.AuthUseCase
}

func NewImageHandler(iuc usecase.ImageUseCase, auc usecase.AuthUseCase) ImageHandler {
	return &imageHandler{
		iuc: iuc,
		auc: auc,
	}
}

type ImageCreateRequest struct {
	SpotID uuid.UUID `json:"spotID"`
	URL    string    `json:"url"`
}

type ImageGetResponse struct {
	Images []model.Image `json:"images"`
}

func (ih *imageHandler) HandleImageGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotID := r.URL.Query().Get("spot_id")

	imgs, err := ih.iuc.GetSpotImgURLBySpotID(ctx, spotID)
	if err != nil {
		http.Error(w, "Failed to get images by spot id", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ImageGetResponse{Images: imgs}); err != nil {
		http.Error(w, "Failed to encode spots to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ih *imageHandler) HandleImageCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ih.auc.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody ImageCreateRequest
	if ok := isValidateImageCreateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid image create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = ih.iuc.ImageCreate(ctx, requestBody.SpotID, requestBody.URL, *user); err != nil {
		http.Error(w, "Internal server error while creating image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateImageCreateRequest(body io.ReadCloser, requestBody *ImageCreateRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.SpotID.String() == DefaultUUID || requestBody.URL == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ih *imageHandler) HandleImageDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ih.auc.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	ok, id, userID := isValidateImageDeleteRequest(r)
	if !ok {
		http.Error(w, "Invalid image delete request", http.StatusBadRequest)
		return
	}

	if err = ih.iuc.ImageDelete(ctx, id, userID, *user); err != nil {
		http.Error(w, "Internal server error while deleting image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateImageDeleteRequest(r *http.Request) (bool, string, string) {
	id := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")

	if id == "" || userID == "" {
		log.Printf("Missing required fields")
		return false, "", ""
	}
	return true, id, userID
}
