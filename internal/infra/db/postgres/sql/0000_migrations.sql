CREATE TABLE IF NOT EXISTS schema_migrations
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);