package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ariefsibuea/articles-feed/internal/api/domain"
	"github.com/ariefsibuea/articles-feed/internal/api/repository"
	_errors "github.com/ariefsibuea/articles-feed/internal/pkg/errors"
)

type ArticleUseCase struct {
	articleRepository repository.ArticleRepository
	authorRepository  repository.AuthorRepository
}

func InitArticleUseCase(articleRepository repository.ArticleRepository, authorRepository repository.AuthorRepository) ArticleUseCase {
	return ArticleUseCase{
		articleRepository: articleRepository,
		authorRepository:  authorRepository,
	}
}

func (u *ArticleUseCase) Create(ctx context.Context, article domain.Article) (domain.Article, error) {
	author, err := u.authorRepository.GetByName(ctx, article.AuthorName)
	if err != nil && !errors.Is(err, _errors.ErrAuthorNotFound) {
		return domain.Article{}, err
	}

	article.AuthorUUID = author.UUID

	if author.UUID == "" {
		newAuthor := domain.Author{
			Name: article.AuthorName,
		}

		authorUUID, err := u.authorRepository.Create(ctx, newAuthor)
		if err != nil {
			return domain.Article{}, err
		}

		article.AuthorUUID = authorUUID
	}

	article.CreatedAt = time.Now().In(time.UTC)
	newArticleUUID, err := u.articleRepository.Create(ctx, article)
	if err != nil {
		return domain.Article{}, err
	}

	article.UUID = newArticleUUID
	return article, nil
}

func (u *ArticleUseCase) GetArticles(ctx context.Context, filter domain.ArticleFilter) (domain.ArticleList, error) {
	return u.articleRepository.GetArticles(ctx, filter)
}
