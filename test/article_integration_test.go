package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ariefsibuea/articles-feed/internal/api/handler"
	"github.com/ariefsibuea/articles-feed/internal/api/repository"
	"github.com/ariefsibuea/articles-feed/internal/api/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	_suite "github.com/stretchr/testify/suite"
)

func (suite *ArticlesFeedTestSuite) SetupTest() {
	suite.cleanupData()

	e := echo.New()

	articleRepository := repository.InitArticleRepository(suite.dbpool)
	authorRepository := repository.InitAuthorRepository(suite.dbpool)

	articleUseCase := usecase.InitArticleUseCase(articleRepository, authorRepository)

	handler.InitArticleHandler(e, articleUseCase)

	suite.echo = e
}

func TestArticlesFeed(t *testing.T) {
	_suite.Run(t, new(ArticlesFeedTestSuite))
}

func (suite *ArticlesFeedTestSuite) TestCreateArticle_Success() {
	payload := map[string]interface{}{
		"title":      "Async Programming in Go",
		"body":       "Understanding goroutines and channels.",
		"authorName": "Evelyn Parker",
	}

	payloadBytes, err := json.Marshal(payload)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewReader(payloadBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	suite.echo.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var createdArticle map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &createdArticle)
	suite.Require().NoError(err)

	isSuccess, _ := createdArticle["success"].(bool)
	assert.True(suite.T(), isSuccess)
	assert.NotNil(suite.T(), createdArticle["data"])

	createdData, _ := createdArticle["data"].(map[string]interface{})
	assert.NotEmpty(suite.T(), createdData["id"])
}

func (suite *ArticlesFeedTestSuite) TestGetArticles_Success() {
	suite.seedArticlesAndAuthors()

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	suite.echo.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var getResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &getResponse)
	suite.Require().NoError(err)

	isSuccess, _ := getResponse["success"].(bool)
	assert.True(suite.T(), isSuccess)
	assert.NotNil(suite.T(), getResponse["data"])
	assert.NotNil(suite.T(), getResponse["meta"])

	meta, _ := getResponse["meta"].(map[string]interface{})
	totalItems, _ := meta["totalItems"].(float64)
	assert.Equal(suite.T(), float64(4), totalItems)

	data, _ := getResponse["data"].(map[string]interface{})
	assert.NotNil(suite.T(), data["articles"])
	articles, _ := data["articles"].([]interface{})
	assert.Equal(suite.T(), 4, len(articles))
}

func (suite *ArticlesFeedTestSuite) TestGetArticlesWithPage_Success() {
	suite.seedArticlesAndAuthors()

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	queryParam := req.URL.Query()
	queryParam.Add("page", "1")
	queryParam.Add("pageSize", "2")
	req.URL.RawQuery = queryParam.Encode()

	suite.echo.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var getResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &getResponse)
	suite.Require().NoError(err)

	isSuccess, _ := getResponse["success"].(bool)
	assert.True(suite.T(), isSuccess)
	assert.NotNil(suite.T(), getResponse["data"])
	assert.NotNil(suite.T(), getResponse["meta"])

	meta, _ := getResponse["meta"].(map[string]interface{})
	totalItems, _ := meta["totalItems"].(float64)
	assert.Equal(suite.T(), float64(4), totalItems)

	data, _ := getResponse["data"].(map[string]interface{})
	assert.NotNil(suite.T(), data["articles"])
	articles, _ := data["articles"].([]interface{})
	assert.Equal(suite.T(), 2, len(articles))
}

func (suite *ArticlesFeedTestSuite) TestGetArticlesWithSearch_Success() {
	suite.seedArticlesAndAuthors()

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	queryParam := req.URL.Query()
	queryParam.Add("query", "PostgreSQL")
	req.URL.RawQuery = queryParam.Encode()

	suite.echo.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var getResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &getResponse)
	suite.Require().NoError(err)

	isSuccess, _ := getResponse["success"].(bool)
	assert.True(suite.T(), isSuccess)
	assert.NotNil(suite.T(), getResponse["data"])
	assert.NotNil(suite.T(), getResponse["meta"])

	meta, _ := getResponse["meta"].(map[string]interface{})
	totalItems, _ := meta["totalItems"].(float64)
	assert.Equal(suite.T(), float64(1), totalItems)

	data, _ := getResponse["data"].(map[string]interface{})
	assert.NotNil(suite.T(), data["articles"])
	articles, _ := data["articles"].([]interface{})
	assert.Equal(suite.T(), 1, len(articles))
}

func (suite *ArticlesFeedTestSuite) seedArticlesAndAuthors() {
	authors := map[string]uuid.UUID{
		"Alice Smith": uuid.New(),
		"Bob Johnson": uuid.New(),
		"Charlie Lee": uuid.New(),
		"Dana White":  uuid.New(),
	}

	query := `
		INSERT INTO authors (author_uuid, name)
		VALUES
			($1, $2),
			($3, $4),
			($5, $6),
			($7, $8)
	`

	args := []interface{}{
		authors["Alice Smith"], "Alice Smith",
		authors["Bob Johnson"], "Bob Johnson",
		authors["Charlie Lee"], "Charlie Lee",
		authors["Dana White"], "Dana White",
	}

	_, err := suite.dbpool.Exec(suite.ctx, query, args...)
	suite.Require().NoError(err)

	query = `
		INSERT INTO articles (article_uuid, author_uuid, title, body, created_at)
		VALUES
			($1, $2, $3, $4, $5),
			($6, $7, $8, $9, $10),
			($11, $12, $13, $14, $15),
			($16, $17, $18, $19, $20)
	`

	args = []interface{}{
		uuid.New(), authors["Alice Smith"], "Introduction to Go", "A quick start guide to Go.", time.Now().UTC(),
		uuid.New(), authors["Bob Johnson"], "Understanding REST APIs", "Learn the basics of RESTful services.", time.Now().UTC(),
		uuid.New(), authors["Charlie Lee"], "Testing in Go", "How to write unit and integration tests.", time.Now().UTC(),
		uuid.New(), authors["Dana White"], "Working with PostgreSQL", "Connecting Go with PostgreSQL.", time.Now().UTC(),
	}

	_, err = suite.dbpool.Exec(suite.ctx, query, args...)
	suite.Require().NoError(err)
}

func (suite *ArticlesFeedTestSuite) cleanupData() {
	_, err := suite.dbpool.Exec(suite.ctx, "TRUNCATE TABLE articles RESTART IDENTITY")
	suite.Require().NoError(err)

	_, err = suite.dbpool.Exec(suite.ctx, "TRUNCATE TABLE authors RESTART IDENTITY")
	suite.Require().NoError(err)
}
