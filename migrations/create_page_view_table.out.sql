CREATE TABLE IF NOT EXISTS event (
  id serial PRIMARY KEY,
  page varchar not null,
  time_spent int not null,
  view_start date not null,
  session_id int not null,
  constraint fk_session_to_pageview FOREIGN KEY(session_id) REFERENCES session(id)
)
