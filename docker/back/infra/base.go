package infra

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	// Register MySQL dialect for goqu
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"

	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type base[T any] struct {
	db        *sql.DB
	dialect   *goqu.DialectWrapper
	tableName string
}

func newBase[T any](db *sql.DB, dialect *goqu.DialectWrapper, tableName string) repository.Base[T] {
	return &base[T]{
		db:        db,
		dialect:   dialect,
		tableName: tableName,
	}
}

func (b *base[T]) List(ctx context.Context, qcs []repository.QueryCondition) ([]T, error) {
	var entitys []T
	var whereClauses []goqu.Expression
	for _, qc := range qcs {
		whereClauses = append(whereClauses, goqu.C(qc.Field).Eq(qc.Value))
	}

	query, _, err := b.dialect.From(b.tableName).Select("*").Where(whereClauses...).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, query)
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

func (b *base[T]) Get(ctx context.Context, id string) (T, error) {
	var entity T
	query, _, err := b.dialect.From(b.tableName).Select("*").Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return entity, err
	}
	err = b.db.QueryRowContext(ctx, query).Scan(&entity)
	return entity, err
}

func (b *base[T]) Create(ctx context.Context, entity T) error {
	query, _, err := b.dialect.Insert(b.tableName).Rows(entity).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) BatchCreate(ctx context.Context, entitys []T) error {
	query, _, err := b.dialect.Insert(b.tableName).Rows(entitys).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) Update(ctx context.Context, id string, entity T) error {
	query, _, err := b.dialect.Update(b.tableName).Set(entity).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) Delete(ctx context.Context, id string) error {
	query, _, err := b.dialect.Delete(b.tableName).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}
