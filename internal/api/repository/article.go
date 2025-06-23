package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ariefsibuea/articles-feed/internal/api/domain"
	_errors "github.com/ariefsibuea/articles-feed/internal/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepository struct {
	dbpool *pgxpool.Pool
}

func InitArticleRepository(dbpool *pgxpool.Pool) ArticleRepository {
	return ArticleRepository{
		dbpool: dbpool,
	}
}

func (r *ArticleRepository) Create(ctx context.Context, article domain.Article) (string, error) {
	_, err := r.dbpool.Exec(ctx, "SET search_path to articles_feed, public")
	if err != nil {
		return "", _errors.ErrInvalidSearchPath
	}

	query := "INSERT INTO articles (author_uuid, title, body, created_at) VALUES ($1, $2, $3, $4) RETURNING article_uuid"
	args := []interface{}{
		article.AuthorUUID,
		article.Title,
		article.Body,
		article.CreatedAt,
	}

	article_uuid := ""
	err = r.dbpool.QueryRow(ctx, query, args...).Scan(&article_uuid)
	if err != nil {
		return "", err
	}

	return article_uuid, nil
}

func (r *ArticleRepository) GetArticles(ctx context.Context, filter domain.ArticleFilter) (domain.ArticleList, error) {
	_, err := r.dbpool.Exec(ctx, "SET search_path to articles_feed, public")
	if err != nil {
		return domain.ArticleList{}, _errors.ErrInvalidSearchPath
	}

	argCounter := 1
	args := make([]interface{}, 0)
	whereCondition := make([]string, 0)

	if q := strings.TrimSpace(filter.Query); q != "" {
		whereCondition = append(whereCondition, fmt.Sprintf(
			"to_tsvector ('simple', coalesce(art.title, '') || ' ' || coalesce(art.body, '')) @@ plainto_tsquery('simple', $%d)", argCounter))
		args = append(args, q)
		argCounter++
	}

	if q := strings.TrimSpace(filter.AuthorName); q != "" {
		whereCondition = append(whereCondition, fmt.Sprintf(
			"to_tsvector('simple', aut.name) @@ plainto_tsquery('simple', $%d)", argCounter))
		args = append(args, q)
		argCounter++
	}

	whereClause := ""
	if len(whereCondition) > 0 {
		whereClause += " WHERE " + strings.Join(whereCondition, " AND ")
	}

	countQuery := `SELECT COUNT (art.article_uuid)
		FROM articles art
		LEFT JOIN authors aut ON art.author_uuid = aut.author_uuid` + whereClause

	var totalItems int32
	err = r.dbpool.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return domain.ArticleList{}, err
	}

	query := `SELECT art.article_uuid, art.title, art.body, art.created_at, aut.name
		FROM articles art
		LEFT JOIN authors aut ON art.author_uuid = aut.author_uuid` + whereClause

	// order by
	query += " ORDER BY art.created_at DESC"

	// limit for page size
	query += fmt.Sprintf(" LIMIT $%d", argCounter)
	args = append(args, filter.PageSize)
	argCounter += 1

	// offset for page
	query += fmt.Sprintf(" OFFSET $%d", argCounter)
	args = append(args, filter.PageSize*(filter.Page-1))

	rows, err := r.dbpool.Query(ctx, query, args...)
	if err != nil {
		return domain.ArticleList{}, err
	}
	defer rows.Close()

	articles := make([]domain.Article, 0)
	for rows.Next() {
		article := domain.Article{}
		articleBody := sql.NullString{}
		authorName := ""

		err := rows.Scan(
			&article.UUID,
			&article.Title,
			&articleBody,
			&article.CreatedAt,
			&authorName,
		)
		if err != nil {
			return domain.ArticleList{}, err
		}

		article.Body = articleBody.String
		article.AuthorName = authorName
		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return domain.ArticleList{}, rows.Err()
	}

	return domain.ArticleList{
		Articles:   articles,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: totalItems,
	}, nil
}
