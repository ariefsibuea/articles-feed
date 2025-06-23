package domain

import "time"

type Article struct {
	UUID       string
	AuthorUUID string
	AuthorName string
	Title      string
	Body       string
	CreatedAt  time.Time
}

type ArticleList struct {
	Articles   []Article
	Page       int32
	PageSize   int32
	TotalItems int32
}

type ArticleFilter struct {
	Page       int32
	PageSize   int32
	Query      string
	AuthorName string
}
