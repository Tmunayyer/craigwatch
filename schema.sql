create table if not exists search (
	id serial primary key,
    "name" text,
    "url" text,
    created_on  timestamptz
);

create table if not exists listing (
    id serial primary key,
    search_id INTEGER REFERENCES search(id),
	data_pid text UNIQUE NOT NULL,
    data_repost_of text,
    unix_date bigint,
    title text,
    link text,
    price text,
    hood text
);