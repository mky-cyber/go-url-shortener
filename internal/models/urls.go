package models

import (
	"database/sql"
	"errors"
)

type URLShortener struct {
	URLs map[string]string
}

type ShortenerDataInterface interface {
	Get(shortened string) (*ShortenerData, error)
	IncreaseClicks(shortened string) error
	Insert(original string, shortenedKey string, clicks int) (int, error)
}

type ShortenerData struct {
	OriginalURL     string
	ShortenedURLKEY string
	Clicks          int
}

type ShortenerDBModel struct {
	DB *sql.DB
}

// Get retrieves a record from the urls table identifying that record by the shortened URL
func (m *ShortenerDBModel) Get(shortenedKey string) (*ShortenerData, error) {
	query := `SELECT original_url, shortened_url_key, clicks FROM urls WHERE shortened_url_key = ?`
	row := m.DB.QueryRow(query, shortenedKey)
	data := &ShortenerData{}
	err := row.Scan(&data.OriginalURL, &data.ShortenedURLKEY, &data.Clicks)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cannot find the matching record")
		}
		return nil, err
	}
	return data, nil
}

// IncreaseClicks increase the clicks number of a given key by one
func (m *ShortenerDBModel) IncreaseClicks(shortenedKey string) error {
	stmt := `UPDATE urls SET clicks = clicks + 1 WHERE shortened_url_key = ?`
	_, err := m.DB.Exec(stmt, shortenedKey)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new record into the urls table
func (m *ShortenerDBModel) Insert(original string, shortenedKey string, clicks int) (int, error) {
	stmt := `INSERT INTO urls  (original_url, shortened_url_key, clicks) VALUES(?, ?, ?)`
	result, err := m.DB.Exec(stmt, original, shortenedKey, clicks)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
