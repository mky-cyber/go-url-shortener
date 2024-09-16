package mocks

import (
	"errors"
	"go-url-shortener/internal/models"
)

type MockShortenerData struct {
	MockData map[string]*models.ShortenerData
}

func (m *MockShortenerData) Get(shortened string) (*models.ShortenerData, error) {
	if data, ok := m.MockData[shortened]; ok {
		return data, nil
	}
	return nil, errors.New("shortened URL not found")
}

func (m *MockShortenerData) GetByOriginalURL(originalURL string) (*models.ShortenerData, error) {
	if data, ok := m.MockData[originalURL]; ok {
		return data, nil
	}
	return nil, errors.New("shortened URL not found")
}

func (m *MockShortenerData) IncreaseClicks(shortened string) error {
	if data, ok := m.MockData[shortened]; ok {
		data.Clicks++
		return nil
	}
	return errors.New("shortened URL not found")
}

func (m *MockShortenerData) Insert(originalURL string, clicks int) (string, string, error) {
	switch originalURL {
	case "https://amazon.com/": // a valid case
		return "abcabc1234567890", "URL successfully shortened", nil
	case "https://google.com/": // assume this url is already shoertened
		return "abcabc1234567890", "URL is already shortened", nil
	default:
		return "", "", errors.New("failed to create the shortened URL")
	}
}
