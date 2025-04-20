create table api_keys (
  uuid text primary key,
  customer_id bigint not null unique,
  inserted timestamp(0)
  with
    time zone not null default now ()
);

create index on api_keys (customer_id);
