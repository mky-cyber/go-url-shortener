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

func (m *MockShortenerData) IncreaseClicks(shortened string) error {
	if data, ok := m.MockData[shortened]; ok {
		data.Clicks++
		return nil
	}
	return errors.New("shortened URL not found")
}

func (m *MockShortenerData) Insert(originalURL string, shortened string, clicks int) (int, error) {
	if _, ok := m.MockData[originalURL]; ok {
		return 1, nil
	}
	return 0, errors.New("failed to create the shortened URL")
}
