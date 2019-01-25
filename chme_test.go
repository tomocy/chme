package chme

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi"
)

func TestChangePostToHiddenMethod(t *testing.T) {
	r := chi.NewRouter()
	r.Use(ChangePostToHiddenMethod)
	r.Post("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("post"))
	})
	r.Put("/put", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("put"))
	})
	r.Patch("/patch", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("patch"))
	})
	r.Delete("/delete", func(w http.ResponseWriter, r *http.Request) {
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
			http.MethodPost,
			"/post",
			url.Values{},
			http.StatusOK,
			"post",
		},
		{
			http.MethodPut,
			"/put",
			url.Values{
				"_method": {"PUT"},
			},
			http.StatusOK,
			"put",
		},
		{
			http.MethodPatch,
			"/patch",
			url.Values{
				"_method": {"PATCH"},
			},
			http.StatusOK,
			"patch",
		},
		{
			http.MethodDelete,
			"/delete",
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
		resp, body := testRequest(t, ts, test.method, test.path, nil)
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
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
