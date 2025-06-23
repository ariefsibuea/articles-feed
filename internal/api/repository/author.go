package repository

import (
	"context"

	"github.com/ariefsibuea/articles-feed/internal/api/domain"
	_errors "github.com/ariefsibuea/articles-feed/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepository struct {
	dbpool *pgxpool.Pool
}

func InitAuthorRepository(dbpool *pgxpool.Pool) AuthorRepository {
	return AuthorRepository{
		dbpool: dbpool,
	}
}

func (r *AuthorRepository) Create(ctx context.Context, author domain.Author) (string, error) {
	_, err := r.dbpool.Exec(ctx, "SET search_path to articles_feed, public")
	if err != nil {
		return "", _errors.ErrInvalidSearchPath
	}

	query := "INSERT INTO authors (name) VALUES ($1) RETURNING author_uuid"
	args := []interface{}{author.Name}

	author_uuid := ""
	err = r.dbpool.QueryRow(ctx, query, args...).Scan(&author_uuid)
	if err != nil {
		return "", err
	}

	return author_uuid, nil
}

func (r *AuthorRepository) GetByName(ctx context.Context, name string) (domain.Author, error) {
	_, err := r.dbpool.Exec(ctx, "SET search_path to articles_feed, public")
	if err != nil {
		return domain.Author{}, _errors.ErrInvalidSearchPath
	}

	query := "SELECT author_uuid, name FROM authors WHERE name = $1"
	args := []interface{}{name}

	author := domain.Author{}

	err = r.dbpool.QueryRow(ctx, query, args...).Scan(
		&author.UUID,
		&author.Name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Author{}, _errors.ErrAuthorNotFound
		}
		return domain.Author{}, err
	}

	return author, nil
}
