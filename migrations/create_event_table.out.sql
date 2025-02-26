CREATE TABLE IF NOT EXISTS pageview (
  id serial primary key,
  data jsonb not null,
  name varchar not null,
  session_id int not null,
  constraint fk_session_to_event FOREIGN KEY(session_id) REFERENCES session(id)
)
