create table config(
  id serial primary key,
  config_id varchar not null,
  data bytea not null,
  data_string varchar,
  env varchar not null default '',
  unique(config_id, env)
);
