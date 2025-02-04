CREATE TABLE IF NOT EXISTS blocked_domains (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    recursive BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

---- create above / drop below ----

DROP TABLE IF EXISTS blocked_domains;

