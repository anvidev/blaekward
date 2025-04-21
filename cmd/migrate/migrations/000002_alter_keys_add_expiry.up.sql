alter table if exists api_keys
add column expiry timestamp(0)
with
  time zone not null default now ();
