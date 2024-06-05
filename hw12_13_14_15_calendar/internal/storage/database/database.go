package database

import (
	"context"
	"fmt"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	conf *config.DatabaseConf
	pg   *pgxpool.Pool
}

type IDatabase interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	DB() *pgxpool.Pool
}

func New(conf *config.DatabaseConf) *DB {
	return &DB{
		conf: conf,
	}
}

func (d *DB) DB() *pgxpool.Pool {
	return d.pg
}

func (d *DB) Connect(ctx context.Context) error {
	conn, err := pgxpool.New(ctx, d.conf.GetDsn())
	if err != nil {
		return fmt.Errorf("database open connect: %w", err)
	}

	d.pg = conn

	return nil
}

func (d *DB) Close(_ context.Context) error {
	d.pg.Close()

	return nil
}
