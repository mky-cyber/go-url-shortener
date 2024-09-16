package models

import (
	"database/sql"
	"errors"
	"go-url-shortener/internal/utils"
	"strings"
	"time"
)

type URLShortener struct {
	URLs map[string]string
}

type ShortenerDataInterface interface {
	Get(shortened string) (*ShortenerData, error)
	GetByOriginalURL(originalURL string) (*ShortenerData, error)
	IncreaseClicks(shortened string) error
	Insert(original string, clicks int) (string, string, error)
}

type ShortenerData struct {
	OriginalURL     string
	ShortenedURLKEY string
	Clicks          int
}

type ShortenerDBModel struct {
	DB *sql.DB
}

const MaxRetry = 5

// Get retrieves a record from the urls table identifying that record by the shortened URL
func (m *ShortenerDBModel) Get(shortenedKey string) (*ShortenerData, error) {
	query := `SELECT original_url, shortened_url_key, clicks FROM urls WHERE shortened_url_key = ?`
	row := m.DB.QueryRow(query, shortenedKey)
	return get(row)
}

func (m *ShortenerDBModel) GetByOriginalURL(originalURL string) (*ShortenerData, error) {
	query := `SELECT original_url, shortened_url_key, clicks FROM urls WHERE original_url = ?`
	row := m.DB.QueryRow(query, originalURL)
	return get(row)
}

func get(r *sql.Row) (*ShortenerData, error) {
	data := &ShortenerData{}
	err := r.Scan(&data.OriginalURL, &data.ShortenedURLKEY, &data.Clicks)
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
	query := `UPDATE urls SET clicks = clicks + 1 WHERE shortened_url_key = ?`
	_, err := m.DB.Exec(query, shortenedKey)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new record into the urls table
// Need to returns 3 arguments shortenedURLKey, responseMessage and error to handle cases like
// case 1: original url is already shortened, return the shortened url key
// case 2: generated key is already used for another url, retry the key generation for a max 5 times
func (m *ShortenerDBModel) Insert(originalURL string, clicks int) (string, string, error) {
	var shortenedKey string
	query := `INSERT INTO urls  (original_url, shortened_url_key, clicks) VALUES(?, ?, ?)`
	// retry for max 5 times to avoid same shortened key though the chance of that happening
	// is very low as we use 16-digits number and letter combinations
	for i := 0; i < MaxRetry; i++ {
		// generate a unique key and save it in db
		shortenedKey = utils.GenerateShortURLKey()
		_, err := m.DB.Exec(query, originalURL, shortenedKey, clicks)
		if err != nil {
			// TODO find a better way to handle duplicate keys
			if strings.Contains(err.Error(), "UNIQUE constraint failed: urls.original_url") {
				data, _ := m.GetByOriginalURL(originalURL)
				return data.ShortenedURLKEY, "URL is already shortened", nil
			}
			if strings.Contains(err.Error(), "UNIQUE constraint failed: urls.shortened_url_key") {
				// small time delay to avoid tight loop
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return "", "", err
		}
	}

	return shortenedKey, "URL successfully shortened", nil
}
