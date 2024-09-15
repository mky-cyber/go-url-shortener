package api

import (
	"go-url-shortener/internal/models"
	"go-url-shortener/internal/models/mocks"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
)

type testServer struct {
	*httptest.Server
}

type testCases struct {
	name                    string
	method                  string
	urlPath                 string
	body                    io.Reader
	expectedStatusCode      int
	expectedResponseMessage string
}

func TestPingRoute(t *testing.T) {
	mockDB := mockDB()
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
				Clicks:          10,
			},
			"https://amazon.com/": {
				OriginalURL:     "https://github.com/",
				ShortenedURLKEY: "abcabc1234567890",
				Clicks:          10,
			},
		},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func runTestCase(
	t *testing.T,
	ts *testServer,
	tc testCases,
) {
	req, err := http.NewRequest(tc.method, ts.URL+tc.urlPath, tc.body)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != tc.expectedStatusCode {
		t.Errorf("got %d; want %d", rs.StatusCode, tc.expectedStatusCode)
	}

	bodyBytes, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	responseMessage := string(bodyBytes)
	if !strings.Contains(responseMessage, tc.expectedResponseMessage) {
		t.Errorf("got %s; want %s", responseMessage, tc.expectedResponseMessage)
	}
}

func TestRedirect(t *testing.T) {
	mockDB := mockDB()
	app := NewApp(mockDB)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	testCases := []testCases{
		{
			name:                    "Valid Redirect",
			method:                  "GET",
			urlPath:                 "/s/abcabc1234567890",
			body:                    nil,
			expectedStatusCode:      http.StatusSeeOther,
			expectedResponseMessage: `<a href="https://github.com/">See Other</a>.`,
		},
		{
			name:                    "URL is invalid",
			method:                  "GET",
			urlPath:                 "/s/invalid-url",
			body:                    nil,
			expectedStatusCode:      http.StatusNotFound,
			expectedResponseMessage: `Shortened URL is invalid`,
		},
		{
			name:                    "URL not found",
			method:                  "GET",
			urlPath:                 "/s/abcabc1234567999",
			body:                    nil,
			expectedStatusCode:      http.StatusNotFound,
			expectedResponseMessage: `Shortened URL not found`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, ts, tc)
		})
	}
}

func TestShortenURL(t *testing.T) {
	mockDB := mockDB()
	app := NewApp(mockDB)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	testCases := []testCases{
		{
			name:                    "Shorten the URL successfully",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"url": "https://amazon.com/"}`),
			expectedStatusCode:      http.StatusOK,
			expectedResponseMessage: "URL successfully shortened",
		},
		{
			name:                    "Invalid JSON payload",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"https://google.com"}`),
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "Invalid JSON payload",
		},
		{
			name:                    "URL is missing",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"url": ""}`),
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "Missing url in the request payload",
		},
		{
			name:                    "Invalid URL",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"url": "htt://google.com"}`),
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "Invalid URL",
		},
		{
			name:                    "Exceeds the maxium length",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"url": "https://www.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com"}`),
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "URL exceeds the maximum length of 2048 characters",
		},
		{
			name:                    "Not reachable URL",
			method:                  "POST",
			urlPath:                 "/shorten",
			body:                    strings.NewReader(`{"url": "https://www.aurlthatprobabilynotexist.com/"}`),
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "The URL was not reachable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, ts, tc)
		})
	}
}
