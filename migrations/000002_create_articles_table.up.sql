create extension if not exists "uuid-ossp";

set search_path = articles_feed, public;

create table if not exists articles (
	id bigserial primary key,
	article_uuid uuid unique not null default uuid_generate_v4(),
    author_uuid uuid,
	title text not null,
	body text,
	created_at timestamp with time zone default now()
);
