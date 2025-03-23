package postgres

import (
	"context"
	"fmt"
	"log"
	"post_service/config"
	"post_service/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db         *pgxpool.Pool
	post       storage.PostRepoI
	attachment storage.PostAttachmentRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (l *Store) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

// Post implements storage.StorageI.
func (s *Store) Post() storage.PostRepoI {
	if s.post == nil {
		s.post = NewPostRepo(s.db)
	}

	return s.post
}

// Session implements storage.StorageI.
func (s *Store) PostAttachment() storage.PostAttachmentRepoI {
	if s.attachment == nil {
		s.attachment = NewAttachmentRepo(s.db)
	}

	return s.attachment
}
