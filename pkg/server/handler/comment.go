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

type CommentCreateRequest struct {
	SpotID   uuid.UUID `json:"spotID"`
	StarRate float64   `json:"starRate"`
	Text     string    `json:"text"`
}

type CommentUpdateRequest struct {
	ID       uuid.UUID `json:"id"`
	SpotID   uuid.UUID `json:"spotID"`
	UserID   uuid.UUID `json:"userID"`
	StarRate float64   `json:"starRate"`
	Text     string    `json:"text"`
}

type CommentDeleteRequest struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userID"`
}

type CommentGetResponse struct {
	Comments []db.Comment
}

type CommentHandler interface {
	HandleCommentCreate(w http.ResponseWriter, r *http.Request)
	HandleCommentGet(w http.ResponseWriter, r *http.Request)
	HandleCommentUpdate(w http.ResponseWriter, r *http.Request)
	HandleCommentDelete(w http.ResponseWriter, r *http.Request)
}

type commentHandler struct {
	cr db.CommentRepository
	ah auth.AuthHandler
}

func NewCommentHandler(cr db.CommentRepository, ah auth.AuthHandler) CommentHandler {
	return &commentHandler{
		cr: cr,
		ah: ah,
	}
}

func (ch *commentHandler) HandleCommentGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	spotID := r.URL.Query().Get("spot_id")

	comments, err := ch.cr.GetCommentBySpotID(ctx, spotID)
	if err != nil {
		http.Error(w, "Failed to get comments by spot id", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CommentGetResponse{Comments: comments}); err != nil {
		http.Error(w, "Failed to encode spots to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ch *commentHandler) HandleCommentCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ユーザー情報取得
	user, err := ch.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody CommentCreateRequest
	if ok := isValidateCommentCreateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := ch.CommentCreate(ctx, user, requestBody); err != nil {
		http.Error(w, "Internal server error while creating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateCommentCreateRequest(body io.ReadCloser, requestBody *CommentCreateRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.SpotID.String() == "" || requestBody.StarRate == 0 || requestBody.Text == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ch *commentHandler) CommentCreate(ctx context.Context, user db.User, requestBody CommentCreateRequest) error {

	var comment = db.Comment{
		SpotID:   requestBody.SpotID,
		UserID:   user.ID,
		StarRate: requestBody.StarRate,
		Text:     requestBody.Text,
	}

	if err := ch.cr.Create(ctx, comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
		return err
	}
	return nil
}

func (ch *commentHandler) HandleCommentUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := ch.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody CommentUpdateRequest
	if ok := isValidateCommentUpdateRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment update request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if user.ID != requestBody.UserID {
		http.Error(w, "User ID mismatch", http.StatusBadRequest)
		return
	}

	if err := ch.CommentUpdate(ctx, requestBody); err != nil {
		http.Error(w, "Internal server error while updating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateCommentUpdateRequest(body io.ReadCloser, requestBody *CommentUpdateRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.ID.String() == "" || requestBody.SpotID.String() == "" || requestBody.UserID.String() == "" || requestBody.StarRate == 0 || requestBody.Text == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}

func (ch *commentHandler) CommentUpdate(ctx context.Context, requestBody CommentUpdateRequest) error {
	var comment = db.Comment{
		ID:       requestBody.ID,
		SpotID:   requestBody.SpotID,
		UserID:   requestBody.UserID,
		StarRate: requestBody.StarRate,
		Text:     requestBody.Text,
	}
	if err := ch.cr.Update(ctx, comment); err != nil {
		log.Printf("Failed to update comment: %v", err)
		return err
	}
	return nil
}

func (ch *commentHandler) HandleCommentDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := ch.ah.FetchUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody CommentDeleteRequest
	if ok := isValidateCommentDeleteRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment delete request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if user.ID != requestBody.UserID {
		http.Error(w, "User ID mismatch", http.StatusBadRequest)
		return
	}

	if err := ch.cr.Delete(ctx, requestBody.ID.String()); err != nil {
		http.Error(w, "Internal server error while deleting comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateCommentDeleteRequest(body io.ReadCloser, requestBody *CommentDeleteRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	if requestBody.ID.String() == "" || requestBody.UserID.String() == "" {
		log.Printf("Missing required fields")
		return false
	}
	return true
}
