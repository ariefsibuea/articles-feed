set search_path = articles_feed, public;

create index if not exists idx_articles_fulltext_search on articles using GIN (
    to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(body, ''))
);
