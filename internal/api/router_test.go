package api

import (
	"go-url-shortener/internal/models"
	"go-url-shortener/internal/models/mocks"
	"go-url-shortener/internal/utils/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPingRoute(t *testing.T) {
	mockDB := mocks.MockDB()
	app := NewApp(mockDB)
	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("got %d; want %d", rs.StatusCode, http.StatusOK)
	}
}

func mockDB() *mocks.MockShortenerData {
	return &mocks.MockShortenerData{
		MockData: map[string]*models.ShortenerData{
			"abcabc1234567890": {
				OriginalURL:     "https://github.com/",
				ShortenedURLKEY: "abcabc1234567890",
				Clicks:          19,
			},
			"https://amazon.com/": {
				OriginalURL:     "https://amazon.com/",
				ShortenedURLKEY: "abcabc1234567890",
				Clicks:          53,
			},
			"https://google.com/": {
				OriginalURL:     "https://google.com/",
				ShortenedURLKEY: "abcabc1234568789",
				Clicks:          10,
			},
		},
	}
}

func TestRedirect(t *testing.T) {
	mockDB := mockDB()
	app := NewApp(mockDB)
	ts := test.NewTestServer(t, app.Routes())
	defer ts.Close()

	testCases := []test.TestCases{
		{
			Name:                    "Valid Redirect",
			Method:                  "GET",
			URLPath:                 "/s/abcabc1234567890",
			Body:                    nil,
			ExpectedStatusCode:      http.StatusSeeOther,
			ExpectedResponseMessage: `<a href="https://github.com/">See Other</a>.`,
		},
		{
			Name:                    "URL is invalid",
			Method:                  "GET",
			URLPath:                 "/s/invalid-url",
			Body:                    nil,
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: `Shortened URL is invalid`,
		},
		{
			Name:                    "URL not found",
			Method:                  "GET",
			URLPath:                 "/s/abcabc1234567999",
			Body:                    nil,
			ExpectedStatusCode:      http.StatusNotFound,
			ExpectedResponseMessage: `Shortened URL not found`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			test.RunTestCase(t, ts, tc)
		})
	}
}

func TestShortenURL(t *testing.T) {
	mockDB := mockDB()
	app := NewApp(mockDB)
	ts := test.NewTestServer(t, app.Routes())
	defer ts.Close()

	testCases := []test.TestCases{
		{
			Name:                    "Shorten the URL successfully",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": "https://amazon.com/"}`),
			ExpectedStatusCode:      http.StatusOK,
			ExpectedResponseMessage: "URL successfully shortened",
		},
		{
			Name:                    "Invalid JSON payload",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"https://google.com"}`),
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: "Invalid JSON payload",
		},
		{
			Name:                    "URL is missing",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": ""}`),
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: "Missing url in the request payload",
		},
		{
			Name:                    "Invalid URL",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": "htt://google.com"}`),
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: "Invalid URL",
		},
		{
			Name:                    "Exceeds the maxium length",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": "https://www.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com"}`),
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: "URL exceeds the maximum length of 2048 characters",
		},
		{
			Name:                    "Not reachable URL",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": "https://www.aurlthatprobabilynotexist.com/"}`),
			ExpectedStatusCode:      http.StatusBadRequest,
			ExpectedResponseMessage: "The URL was not reachable",
		},
		{
			Name:                    "URL already exists",
			Method:                  "POST",
			URLPath:                 "/shorten",
			Body:                    strings.NewReader(`{"url": "https://google.com/"}`),
			ExpectedStatusCode:      http.StatusOK,
			ExpectedResponseMessage: "URL is already shortened",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			test.RunTestCase(t, ts, tc)
		})
	}
}
