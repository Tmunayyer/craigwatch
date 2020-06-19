create table if not exists monitor (
	id serial primary key,
    email text,
    "url" text,
    confirmed boolean,
    interval integer,
    created_on  timestamptz,
	polled_on timestamptz
);