create table if not exists monitor (
	id serial primary key,
    "name" text,
    "url" text,
    confirmed boolean,
    interval integer default 60000,
    created_on  timestamptz
);

create table if not exists listing (
    id serial primary key,
    monitor_id INTEGER REFERENCES monitor(id),
	data_pid text,
    data_repost_of text,
    "date" timestamptz,
    title text,
    link text,
    price text,
    hood text
);