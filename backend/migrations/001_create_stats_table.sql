CREATE TABLE IF NOT EXISTS stats_aggregated (
  id SERIAL PRIMARY KEY,
  time TIMESTAMP NOT NULL,
  domain TEXT NOT NULL,
  count INTEGER NOT NULL,
  q_type TEXT NOT NULL DEFAULT 'A',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (time, domain)
);

---- create above / drop below ----

DROP TABLE IF EXISTS stats_aggregated;
