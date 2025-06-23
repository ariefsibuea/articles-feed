package handler

import (
	"net/http"

	"github.com/ariefsibuea/articles-feed/internal/api/usecase"

	"github.com/labstack/echo/v4"
)

type articleHandler struct {
	articleUseCase usecase.ArticleUseCase
}

func InitArticleHandler(e *echo.Echo, articleUseCase usecase.ArticleUseCase) {
	handler := &articleHandler{
		articleUseCase: articleUseCase,
	}

	e.POST("/articles", handler.create)
	e.GET("/articles", handler.get)
}

func (h *articleHandler) create(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(CreateArticleRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	res, err := h.articleUseCase.Create(ctx, req.ToDomain())
	if err != nil {
		return err
	}

	return Success(c, http.StatusCreated, CreateArticleResponseFromDomain(res), nil)
}

func (h *articleHandler) get(c echo.Context) error {
	ctx := c.Request().Context()

	binder := new(echo.DefaultBinder)
	req := new(GetArticlesRequest)
	if err := binder.BindQueryParams(c, req); err != nil {
		return err
	}

	res, err := h.articleUseCase.GetArticles(ctx, req.ToFilterDomain())
	if err != nil {
		return err
	}

	return Success(c, http.StatusOK, GetArticlesResponseFromDomain(res), &Meta{
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalItems: res.TotalItems,
	})
}
