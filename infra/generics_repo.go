package infra

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/tusmasoma/campfinder/domain/repository"
)

type genericRepository[T any] struct {
	db        *sql.DB
	tableName string
}

func NewGenericRepository[T any](db *sql.DB, tableName string) repository.Repository[T] {
	return &genericRepository[T]{
		db:        db,
		tableName: tableName,
	}
}

func (gr *genericRepository[T]) GetByID(ctx context.Context, id string) (T, error) {
	var entity T
	query := `SELECT * FROM ` + gr.tableName + ` WHERE id = ?`
	err := gr.db.QueryRowContext(ctx, query, id).Scan(&entity)
	return entity, err
}

func (gr *genericRepository[T]) Create(ctx context.Context, entity T) error {
	// リフレクションを使用してエンティティのフィールド名と値を取得
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	var fieldNames []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldNames = append(fieldNames, field.Name)
		placeholders = append(placeholders, "?")
		values = append(values, v.Field(i).Interface())
	}

	// SQLクエリ構築
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		gr.tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := gr.db.ExecContext(ctx, query, values...)

	return err
}

func (gr *genericRepository[T]) Update(ctx context.Context, entity T) error {
	// リフレクションを使用してエンティティのフィールド名と値を取得
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	var fieldNames []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "ID" {
			fieldNames = append(fieldNames, field.Name+"=?")
			values = append(values, v.Field(i).Interface())
		}
	}

	// idフィールドの値を取得 (エンティティにIDフィールドが存在することを前提としています)
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return fmt.Errorf("ID field not found in entity")
	}
	id := idField.Interface()

	// SQLクエリ構築
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		gr.tableName,
		strings.Join(fieldNames, ", "),
	)

	// idの値をvaluesの末尾に追加
	values = append(values, id)

	_, err := gr.db.ExecContext(ctx, query, values...)

	return err
}

func (gr *genericRepository[T]) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM ` + gr.tableName + ` WHERE id = ?`
	_, err := gr.db.ExecContext(ctx, query, id)
	return err
}
