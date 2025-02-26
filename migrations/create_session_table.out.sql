CREATE TABLE IF NOT EXISTS session (
  id serial primary key,
  anonymous_id int,
  timestamp date,
  referrer varchar,
  screen_width int,
  ip varchar,
  user_agent varchar,
  country varchar,
  country_code varchar,
  os varchar,
  browser varchar
)
