package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg DB) (*pgxpool.Pool, error) {
	dsn := dsnFromDBConfig(cfg)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	return pool, nil
}

func dsnFromDBConfig(cfg DB) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

type Queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func FetchRow[Object any](ctx context.Context, q Queryer, sql string, args ...any) (*Object, error) {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	object, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[Object])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrObjectNotFound
	}

	return object, err
}

// FetchRows may come in handy later.
func FetchRows[
	Object any,
](ctx context.Context, q Queryer, sql string, args ...any) ([]*Object, error) {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	objects, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[Object])

	return objects, err
}

func Execute(ctx context.Context, q Queryer, sql string, args ...any) error {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.CommandTag().RowsAffected() > 0 {
		return ErrObjectNotFound
	}

	rows.Next()

	return nil
}
