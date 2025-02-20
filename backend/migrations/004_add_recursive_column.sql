ALTER TABLE stats_aggregated ADD COLUMN q_type VARCHAR(10) NOT NULL DEFAULT '';

---- create above / drop below ----

ALTER TABLE stats_aggregated DROP COLUMN q_type;
