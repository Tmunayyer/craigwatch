create table if not exists monitor (
	id serial primary key,
    email text,
    "url" text,
    confirmed boolean,
    interval integer default 60000,
    created_on  timestamptz
);