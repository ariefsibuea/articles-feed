set search_path = articles_feed, public;

create index if not exists idx_author_name_fulltext_search on authors using GIN (
    to_tsvector('simple', name)
)
