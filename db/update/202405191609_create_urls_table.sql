-- migrate:up
-- Create a urls table that stores the original urls and shortented urls along with the clicks
CREATE TABLE IF NOT EXISTS "urls" (
    url_id INTEGER  PRIMARY KEY AUTOINCREMENT,
    original_url TEXT NOT NULL,
    shortened_url_key TEXT NOT NULL,
    clicks INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_original_url UNIQUE (original_url)
);

CREATE index idx_shortened ON urls (shortened_url_key);

-- DATA
INSERT INTO urls (original_url, shortened_url_key, clicks) 
VALUES
('https://github.com', 'K8FN2X9z48U8koUy', 100),
('https://stackoverflow.com', 'xl1e3PgCs9p6Zde6', 500);
-- migrate:down
