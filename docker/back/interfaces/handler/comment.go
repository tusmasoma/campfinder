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

type CommentHandler interface {
	ListComments(w http.ResponseWriter, r *http.Request)
	CreateComment(w http.ResponseWriter, r *http.Request)
	UpdateComment(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
}

type commentHandler struct {
	cuc usecase.CommentUseCase
	auc usecase.AuthUseCase
}

func NewCommentHandler(cuc usecase.CommentUseCase, auc usecase.AuthUseCase) CommentHandler {
	return &commentHandler{
		cuc: cuc,
		auc: auc,
	}
}

type CreateCommentRequest struct {
	SpotID   uuid.UUID `json:"spotID"`
	StarRate float64   `json:"starRate"`
	Text     string    `json:"text"`
}

type UpdateCommentRequest struct {
	ID       uuid.UUID `json:"id"`
	SpotID   uuid.UUID `json:"spotID"`
	UserID   uuid.UUID `json:"userID"`
	StarRate float64   `json:"starRate"`
	Text     string    `json:"text"`
}

type ListCommentResponse struct {
	Comments []model.Comment `json:"comments"`
}

func (ch *commentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotID := r.URL.Query().Get("spot_id")

	comments, err := ch.cuc.ListComments(ctx, spotID)
	if err != nil {
		http.Error(w, "Failed to get comments by spot id", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ListCommentResponse{Comments: comments}); err != nil {
		http.Error(w, "Failed to encode comments to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ch *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ch.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody CreateCommentRequest
	if ok := isValidateCreateCommentRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = ch.cuc.CreateComment(ctx, requestBody.SpotID, requestBody.StarRate, requestBody.Text, *user); err != nil {
		http.Error(w, "Internal server error while creating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateCreateCommentRequest(body io.ReadCloser, requestBody *CreateCommentRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.SpotID.String() == DefaultUUID || requestBody.StarRate == 0 || requestBody.Text == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ch *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ch.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody UpdateCommentRequest
	if ok := isValidateUpdateCommentRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment update request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = ch.cuc.UpdateComment(
		ctx,
		requestBody.ID,
		requestBody.SpotID,
		requestBody.UserID,
		requestBody.StarRate,
		requestBody.Text,
		*user,
	); err != nil {
		http.Error(w, "Internal server error while updating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateUpdateCommentRequest(body io.ReadCloser, requestBody *UpdateCommentRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.ID.String() == DefaultUUID ||
		requestBody.SpotID.String() == DefaultUUID ||
		requestBody.UserID.String() == DefaultUUID ||
		requestBody.StarRate == 0 ||
		requestBody.Text == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ch *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ch.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	ok, id, userID := isValidateDeleteCommentRequest(r)
	if !ok {
		http.Error(w, "Invalid comment delete request", http.StatusBadRequest)
		return
	}

	if err = ch.cuc.DeleteComment(ctx, id, userID, *user); err != nil {
		http.Error(w, "Internal server error while deleting comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateDeleteCommentRequest(r *http.Request) (bool, string, string) {
	id := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")

	if id == "" || userID == "" {
		log.Printf("Missing required fields")
		return false, "", ""
	}
	return true, id, userID
}
