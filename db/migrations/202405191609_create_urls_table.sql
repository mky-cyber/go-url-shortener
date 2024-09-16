-- migrate:up
-- Create a urls table that stores the original urls and shortented urls along with the clicks
CREATE TABLE IF NOT EXISTS "urls" (
    url_id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_url TEXT UNIQUE NOT NULL,
    shortened_url_key TEXT UNIQUE NOT NULL,
    clicks INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create index on original_url and shortened_url_key as these 2 fields will be quired the most
CREATE INDEX idx_original_url ON urls (original_url);
CREATE index idx_shortened ON urls (shortened_url_key);

-- Create a trigger to update the updated field when a row is modified
CREATE TRIGGER update_timestamp
AFTER UPDATE ON urls
FOR EACH ROW
BEGIN
    UPDATE urls SET updated = CURRENT_TIMESTAMP WHERE url_id = OLD.url_id;
END;

-- DATA
INSERT INTO urls (original_url, shortened_url_key, clicks, active) 
VALUES
('https://github.com', 'K8FN2X9z48U8koUy', 100, TRUE),
('https://stackoverflow.com', 'xl1e3PgCs9p6Zde6', 500, TRUE),
('https://amazon.com', '123e3PgCs9p6Zabc', 0, FALSE); -- set active to false, an example for soft delete
-- migrate:down
