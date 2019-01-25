package chme

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestChangePostToHiddenMethod(t *testing.T) {
	endpoint := "/chme"
	r := chi.NewRouter()
	r.Use(ChangePostToHiddenMethod)
	r.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get"))
	})
	r.Post(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("post"))
	})
	r.Put(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("put"))
	})
	r.Patch(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patch"))
	})
	r.Delete(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("delete"))
	})

	tests := []struct {
		method             string
		path               string
		body               url.Values
		expectedStatusCode int
		expectedBody       string
	}{
		{
			http.MethodGet,
			endpoint,
			url.Values{},
			http.StatusOK,
			"get",
		},
		{
			http.MethodPost,
			endpoint,
			url.Values{},
			http.StatusOK,
			"post",
		},
		{
			http.MethodPost,
			endpoint,
			url.Values{
				"_method": {"PUT"},
			},
			http.StatusOK,
			"put",
		},
		{
			http.MethodPost,
			endpoint,
			url.Values{
				"_method": {"PATCH"},
			},
			http.StatusOK,
			"patch",
		},
		{
			http.MethodPost,
			endpoint,
			url.Values{
				"_method": {"DELETE"},
			},
			http.StatusOK,
			"delete",
		},
	}

	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		resp, body := testRequest(t, ts, test.method, test.path, strings.NewReader(test.body.Encode()))
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("unexpected status code: got %d, but expected %d\n", resp.StatusCode, test.expectedStatusCode)
		}
		if body != test.expectedBody {
			t.Errorf("unexpected body: got %s, but expected %s\n", body, test.expectedBody)
		}
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatalf("failed to test: faield to create new request: %s\n", err)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to test: failed to request: %s\n", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to test: failed to read all of response body: %s\n", err)
	}

	return resp, string(respBody)
}
