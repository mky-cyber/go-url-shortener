package test

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestServer struct {
	*httptest.Server
}

type TestCases struct {
	Name                    string
	Method                  string
	URLPath                 string
	Body                    io.Reader
	ExpectedStatusCode      int
	ExpectedResponseMessage string
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

func RunTestCase(
	t *testing.T,
	ts *TestServer,
	tc TestCases,
) {
	req, err := http.NewRequest(tc.Method, ts.URL+tc.URLPath, tc.Body)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != tc.ExpectedStatusCode {
		t.Errorf("got %d; want %d", rs.StatusCode, tc.ExpectedStatusCode)
	}

	bodyBytes, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	responseMessage := string(bodyBytes)
	if !strings.Contains(responseMessage, tc.ExpectedResponseMessage) {
		t.Errorf("got %s; want %s", responseMessage, tc.ExpectedResponseMessage)
	}
}
