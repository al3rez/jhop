package jhop

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	file = `
	{
		"profile": {
			"name": "foo"
		},
		"recipes": [{
			"difficulty": "hard",
			"id": 1,
			"prep_time": "1h"
		}, {
			"difficulty": "easy",
			"id": 2,
			"prep_time": "15m"
		}]
	}
`
)

func TestNewServer(t *testing.T) {
	type check struct {
		req     *http.Request
		code    int
		content string
	}

	tests := []struct {
		name    string
		file    io.Reader
		wantErr bool
		checks  []check
	}{
		{
			name:    "provide resources",
			file:    strings.NewReader(file),
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
					content: file,
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
			},
		},
		{
			name:    "with error",
			file:    strings.NewReader("///"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewHandler(tt.file)
			if tt.wantErr && err == nil {
				t.Error("expected error")
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
