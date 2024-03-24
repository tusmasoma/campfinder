package infra

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type genericRepository[T any] struct {
	db        *sql.DB
	dialect   *goqu.DialectWrapper
	tableName string
}

func NewGenericRepository[T any](db *sql.DB, dialect *goqu.DialectWrapper, tableName string) repository.Repository[T] {
	return &genericRepository[T]{
		db:        db,
		dialect:   dialect,
		tableName: tableName,
	}
}

func (gr *genericRepository[T]) List(ctx context.Context, qcs []repository.QueryCondition) ([]T, error) {
	var entitys []T
	var whereClauses []goqu.Expression
	for _, qc := range qcs {
		whereClauses = append(whereClauses, goqu.C(qc.Field).Eq(qc.Value))
	}

	query, _, err := gr.dialect.From(gr.tableName).Select("*").Where(whereClauses...).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := gr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity T
		if err = rows.Scan(&entity); err != nil {
			return nil, err
		}
		entitys = append(entitys, entity)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entitys, nil
}

func (gr *genericRepository[T]) Get(ctx context.Context, id string) (T, error) {
	var entity T
	query, _, err := gr.dialect.From(gr.tableName).Select("*").Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return entity, err
	}
	err = gr.db.QueryRowContext(ctx, query).Scan(&entity)
	return entity, err
}

func (gr *genericRepository[T]) Create(ctx context.Context, entitys []T) error {
	query, _, err := gr.dialect.Insert(gr.tableName).Rows(entitys).ToSQL()
	if err != nil {
		return err
	}
	_, err = gr.db.ExecContext(ctx, query)
	return err
}

func (gr *genericRepository[T]) Update(ctx context.Context, id string, entity T) error {
	query, _, err := gr.dialect.Update(gr.tableName).Set(entity).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = gr.db.ExecContext(ctx, query)
	return err
}

func (gr *genericRepository[T]) Delete(ctx context.Context, id string) error {
	query, _, err := gr.dialect.Delete(gr.tableName).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = gr.db.ExecContext(ctx, query)
	return err
}
