create table if not exists search (
	id serial primary key,
    "name" text,
    "url" text unique not null,
    created_on  timestamptz,
    timezone text not null
);

create table if not exists listing (
    id serial primary key,
    search_id integer references search(id),
	data_pid text unique not null,
    data_repost_of text,
    unix_date bigint,
    title text,
    link text,
    price text,
    hood text
);

create or replace function _final_median(numeric[])
   returns numeric as
$$
   select avg(val)
   from (
     select val
     from unnest($1) val
     order by 1
     limit  2 - mod(array_upper($1, 1), 2)
     offset ceil(array_upper($1, 1) / 2.0) - 1
   ) sub;
$$
language 'sql' immutable;

create or replace aggregate median(numeric) (
  sfunc=array_append,
  stype=numeric[],
  finalfunc=_final_median,
  initcond='{}'
);
