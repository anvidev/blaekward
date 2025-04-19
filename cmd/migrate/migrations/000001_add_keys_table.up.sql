create table api_keys (
  uuid text primary key,
  customer_id bigint not null unique,
  inserted timestamptz not null default now ()
);

create index on api_keys (customer_id);
