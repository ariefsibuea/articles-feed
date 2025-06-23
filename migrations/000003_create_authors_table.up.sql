set search_path = articles_feed, public;

create table if not exists authors (
	id bigserial primary key,
	author_uuid uuid unique not null default uuid_generate_v4(),
	name varchar(255) not null
);
