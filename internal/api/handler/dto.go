package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/ariefsibuea/articles-feed/internal/api/domain"
)

type CreateArticleRequest struct {
	Title      string `json:"title"`
	AuthorName string `json:"authorName"`
	Body       string `json:"body"`
}

func (req *CreateArticleRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("'title' is required")
	}
	if strings.TrimSpace(req.AuthorName) == "" {
		return fmt.Errorf("'authorName' is required")
	}
	return nil
}

func (req *CreateArticleRequest) ToDomain() domain.Article {
	return domain.Article{
		AuthorName: req.AuthorName,
		Title:      req.Title,
		Body:       req.Body,
	}
}

type CreateArticleResponse struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"createdAt"`
}

func CreateArticleResponseFromDomain(article domain.Article) CreateArticleResponse {
	return CreateArticleResponse{
		ID:         article.UUID,
		Title:      article.Title,
		AuthorName: article.AuthorName,
		Body:       article.Body,
		CreatedAt:  article.CreatedAt,
	}
}

type GetArticlesRequest struct {
	Page       int32  `query:"page"`
	PageSize   int32  `query:"pageSize"`
	Query      string `query:"query"`
	AuthorName string `query:"authorName"`
}

func (req *GetArticlesRequest) ToFilterDomain() domain.ArticleFilter {
	articleFilter := domain.ArticleFilter{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Query:      req.Query,
		AuthorName: req.AuthorName,
	}

	if articleFilter.Page <= 0 {
		articleFilter.Page = 1
	}
	if articleFilter.PageSize <= 0 {
		articleFilter.PageSize = 20
	}

	return articleFilter
}

type ArticleResponse struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetArticlesResponse struct {
	Articles []ArticleResponse `json:"articles"`
}

func GetArticlesResponseFromDomain(articleList domain.ArticleList) GetArticlesResponse {
	articlesResponse := make([]ArticleResponse, 0, len(articleList.Articles))

	for _, a := range articleList.Articles {
		articleResponse := ArticleResponse{
			ID:         a.UUID,
			AuthorName: a.AuthorName,
			Title:      a.Title,
			Body:       a.Body,
			CreatedAt:  a.CreatedAt,
		}
		articlesResponse = append(articlesResponse, articleResponse)
	}

	return GetArticlesResponse{
		Articles: articlesResponse,
	}
}
