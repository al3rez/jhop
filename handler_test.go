package jhop

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	recipesJSON = `{
		"recipes": [{
			"difficulty": "hard",
			"id": 1,
			"prep_time": "1h"
		}, {
			"difficulty": "easy",
			"id": 2,
			"prep_time": "15m"
		}]
	}`
	profileJSON = `{
		"profile": [{
			"name": "John Doe",
			"id": 1
		}]
	}`
)

func TestNewServer(t *testing.T) {
	type check struct {
		req     *http.Request
		code    int
		content string
	}

	tests := []struct {
		name    string
		files   []io.ReadCloser
		wantErr bool
		checks  []check
	}{
		{
			name: "provide resources",
			files: []io.ReadCloser{
				ioutil.NopCloser(strings.NewReader(recipesJSON)),
				ioutil.NopCloser(strings.NewReader(profileJSON)),
			},
			wantErr: false,
			checks: []check{
				check{
					req:     httptest.NewRequest("GET", "http://localhost/recipes/1", nil),
					code:    200,
					content: "{\"difficulty\":\"hard\",\"id\":1,\"prep_time\":\"1h\"}",
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/recipes/2", nil),
					code:    200,
					content: "{\"difficulty\":\"easy\",\"id\":2,\"prep_time\":\"15m\"}",
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/recipes", nil),
					code:    200,
					content: recipesJSON,
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/recipes/3", nil),
					code:    404,
					content: "",
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/profile/3", nil),
					code:    404,
					content: "",
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/profiles", nil),
					code:    200,
					content: profileJSON,
				},
				check{
					req:     httptest.NewRequest("GET", "http://localhost/profiles/1", nil),
					code:    200,
					content: "{\"name\":\"John Doe\",\"id\":1}",
				},
			},
		},
		{
			name: "with error",
			files: []io.ReadCloser{
				ioutil.NopCloser(strings.NewReader("///")),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routes := map[string]string{
				"/profile/{id}": "/profiles/{id}",
				"/profile":      "/profiles",
			}
			h, err := NewHandlerWithRoutes(routes, tt.files...)
			if tt.wantErr && err == nil {
				t.Error("expected error")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			for _, check := range tt.checks {
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, check.req)

				if check.code != resp.Code {
					t.Errorf("expected code: %d got %d", check.code, resp.Code)
					return
				}

				body := resp.Body.String()
				if strings.Contains(check.content, body) {
					t.Errorf("expected body: %s got %s", check.content, body)
				}
			}
		})
	}
}
