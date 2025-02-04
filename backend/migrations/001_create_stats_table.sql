CREATE TABLE IF NOT EXISTS stats_aggregated (
  id SERIAL PRIMARY KEY,
  time TIMESTAMP NOT NULL,
  domain TEXT NOT NULL,
  count INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (time, domain)
);

---- create above / drop below ----

DROP TABLE IF EXISTS stats_aggregated;
