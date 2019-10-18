create table config(
  id serial primary key,
  config_id varchar not null,
  data jsonb not null default '{}',
  env varchar not null default '',
  unique(config_id, env)
)
