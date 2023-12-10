CREATE TABLE IF NOT EXISTS urls (
    original VARCHAR(2048) UNIQUE NOT NULL,
    alias VARCHAR(128) UNIQUE NOT NULL
);

-- CREATE INDEX idx_origin_fulltext ON urls USING gin (to_tsvector('english', original));
-- CREATE INDEX idx_short_fulltext ON urls USING gin (to_tsvector('english', alias));