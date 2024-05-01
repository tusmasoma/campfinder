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
	BatchCreateComments(w http.ResponseWriter, r *http.Request)
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

type BatchCreateCommentsRequest struct {
	Comments []CreateCommentRequest `json:"comments"`
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

	params := convertCreateCommentReqeuestToParams(requestBody, user.ID)
	if err = ch.cuc.CreateComment(ctx, params); err != nil {
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

func convertCreateCommentReqeuestToParams(req CreateCommentRequest, userID uuid.UUID) *usecase.CreateCommentParams {
	return &usecase.CreateCommentParams{
		UserID:   userID,
		SpotID:   req.SpotID,
		StarRate: req.StarRate,
		Text:     req.Text,
	}
}

func (ch *commentHandler) BatchCreateComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := ch.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	var requestBody BatchCreateCommentsRequest
	if ok := isValidateBatchCreateCommentsRequest(r.Body, &requestBody); !ok {
		http.Error(w, "Invalid comment batch create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := convertBatchCreateCommentsRequestToParams(requestBody, user.ID)
	if err = ch.cuc.BatchCreateComments(ctx, params); err != nil {
		http.Error(w, "Internal server error while batch creating comments", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidateBatchCreateCommentsRequest(body io.ReadCloser, requestBody *BatchCreateCommentsRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		return false
	}
	for _, comment := range requestBody.Comments {
		if comment.SpotID.String() == DefaultUUID || comment.StarRate == 0 || comment.Text == "" {
			log.Printf("Missing required fields")
			return false
		}
	}
	return true
}

func convertBatchCreateCommentsRequestToParams(
	req BatchCreateCommentsRequest,
	userID uuid.UUID,
) *usecase.BatchCreateCommentsParams {
	var comments []usecase.CreateCommentParams
	for _, commentReq := range req.Comments {
		comments = append(comments, usecase.CreateCommentParams{
			UserID:   userID,
			SpotID:   commentReq.SpotID,
			StarRate: commentReq.StarRate,
			Text:     commentReq.Text,
		})
	}
	return &usecase.BatchCreateCommentsParams{
		Comments: comments,
	}
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
