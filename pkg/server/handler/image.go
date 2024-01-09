//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tusmasoma/campfinder/db"
	"github.com/tusmasoma/campfinder/internal/auth"
)

type ImageCreateRequest struct {
	SpotID uuid.UUID `json:"spotID"`
	URL    string    `json:"url"`
}

type ImageDeleteRequest struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userID"`
}

type ImageGetResponse struct {
	Images []db.Image
}

type ImageHandler interface {
	HandleImageGet(w http.ResponseWriter, r *http.Request)
	HandleImageCreate(w http.ResponseWriter, r *http.Request)
	HandleImageDelete(w http.ResponseWriter, r *http.Request)
}

type imageHandler struct {
	ir db.ImageRepository
	ah auth.AuthHandler
}

func NewImageHandler(ir db.ImageRepository, ah auth.AuthHandler) ImageHandler {
	return &imageHandler{
		ir: ir,
		ah: ah,
	}
}

func (ih *imageHandler) HandleImageGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	spotID := r.URL.Query().Get("spot_id")

	imgs, err := ih.ir.GetSpotImgURLBySpotId(ctx, spotID)
	if err != nil {
		http.Error(w, "Failed to get images by spot id", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ImageGetResponse{Images: imgs}); err != nil {
		http.Error(w, "Failed to encode spots to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ih *imageHandler) HandleImageCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ユーザー情報取得
	user, err := ih.ah.FetchUserFromContext(ctx)
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

	if err := ih.ImageCreate(ctx, user, requestBody); err != nil {
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
	if requestBody.SpotID.String() == "00000000-0000-0000-0000-000000000000" || requestBody.URL == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ih *imageHandler) ImageCreate(ctx context.Context, user db.User, requestBody ImageCreateRequest) error {

	var img = db.Image{
		SpotID: requestBody.SpotID,
		UserID: user.ID,
		URL:    requestBody.URL,
	}
	if err := ih.ir.Create(ctx, img); err != nil {
		log.Printf("Failed to create image: %v", err)
		return err
	}
	return nil
}

func (ih *imageHandler) HandleImageDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := ih.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody ImageDeleteRequest
	if ok := isValidateImageDeleteRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid image delete request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !user.IsAdmin && user.ID != requestBody.UserID {
		http.Error(w, "User ID mismatch", http.StatusInternalServerError)
		return
	}

	if err := ih.ir.Delete(ctx, requestBody.ID.String()); err != nil {
		http.Error(w, "Internal server error while deleting image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateImageDeleteRequest(body io.ReadCloser, requestBody *ImageDeleteRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.ID.String() == "00000000-0000-0000-0000-000000000000" || requestBody.UserID.String() == "00000000-0000-0000-0000-000000000000" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}
