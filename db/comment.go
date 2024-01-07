package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CommentRepository interface {
	GetCommentBySpotID(ctx context.Context, spotID uuid.UUID, opts ...QueryOptions) (comments []Comment, err error)
	GetCommentByID(ctx context.Context, id uuid.UUID, opts ...QueryOptions) (comment Comment, err error)
	Create(ctx context.Context, comment Comment, opts ...QueryOptions) (err error)
	Update(ctx context.Context, comment Comment, opts ...QueryOptions) (err error)
	Delete(ctx context.Context, id uuid.UUID, opts ...QueryOptions) (err error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

type Comment struct {
	ID       uuid.UUID
	SpotID   uuid.UUID
	UserID   uuid.UUID
	StarRate float64 `json:"starRate"`
	Text     string  `json:"text"`
	Created  time.Time
}

func (cr *commentRepository) GetCommentBySpotID(ctx context.Context, spotID uuid.UUID, opts ...QueryOptions) (comments []Comment, err error) {
	var executor SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Comment
	WHERE spot_id=?
	`
	rows, err := executor.QueryContext(ctx, query, spotID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.ID,
			&comment.SpotID,
			&comment.UserID,
			&comment.StarRate,
			&comment.Text,
			&comment.Created,
		)
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (cr *commentRepository) GetCommentByID(ctx context.Context, id uuid.UUID, opts ...QueryOptions) (comment Comment, err error) {
	var executor SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Comment
	WHERE id=?
	`
	err = executor.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.SpotID,
		&comment.UserID,
		&comment.StarRate,
		&comment.Text,
		&comment.Created,
	)
	return
}

func (cr *commentRepository) Create(ctx context.Context, comment Comment, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	INSERT INTO Comment (
		spot_id, user_id, star_rate, text
		)
		VALUES (?, ?, ?, ?)
		RETURNING id;
		`
	err = executor.QueryRowContext(
		ctx,
		query,
		comment.SpotID,
		comment.UserID,
		comment.StarRate,
		comment.Text,
	).Scan(&comment.ID)

	return
}

func (cr *commentRepository) Update(ctx context.Context, comment Comment, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	UPDATE Comment SET
	star_rate=?,text=?
	WHERE id=?
	`
	_, err = executor.ExecContext(ctx, query, comment.StarRate, comment.Text, comment.ID)
	return
}

func (cr *commentRepository) Delete(ctx context.Context, id uuid.UUID, opts ...QueryOptions) (err error) {
	var executor SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}
	_, err = executor.ExecContext(ctx, "DELETE FROM Comment WHERE id = ?", id)
	return
}
