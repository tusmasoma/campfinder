//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/tusmasoma/campfinder/docker/back/domain/model"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) repository.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (cr *commentRepository) GetCommentBySpotID(
	ctx context.Context,
	spotID string,
	opts ...repository.QueryOptions,
) ([]model.Comment, error) {
	var executor repository.SQLExecutor = cr.db
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
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err = rows.Scan(
			&comment.ID,
			&comment.SpotID,
			&comment.UserID,
			&comment.StarRate,
			&comment.Text,
			&comment.Created,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) GetCommentByID(
	ctx context.Context,
	id string,
	opts ...repository.QueryOptions,
) (model.Comment, error) {
	var executor repository.SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	SELECT *
	FROM Comment
	WHERE id=?
	`
	var comment model.Comment
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.SpotID,
		&comment.UserID,
		&comment.StarRate,
		&comment.Text,
		&comment.Created,
	)
	return comment, err
}

func (cr *commentRepository) Create(ctx context.Context, comment model.Comment, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	INSERT INTO Comment (
		id, spot_id, user_id, star_rate, text
		)
		VALUES (?, ?, ?, ?, ?)
		`
	_, err := executor.ExecContext(
		ctx,
		query,
		uuid.New(),
		comment.SpotID,
		comment.UserID,
		comment.StarRate,
		comment.Text,
	)

	return err
}

func (cr *commentRepository) Update(ctx context.Context, comment model.Comment, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}

	query := `
	UPDATE Comment SET
	star_rate=?,text=?
	WHERE id=?
	`
	_, err := executor.ExecContext(ctx, query, comment.StarRate, comment.Text, comment.ID)
	return err
}

func (cr *commentRepository) Delete(ctx context.Context, id string, opts ...repository.QueryOptions) error {
	var executor repository.SQLExecutor = cr.db
	if len(opts) > 0 && opts[0].Executor != nil {
		executor = opts[0].Executor
	}
	_, err := executor.ExecContext(ctx, "DELETE FROM Comment WHERE id = ?", id)
	return err
}
