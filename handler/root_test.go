package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "root",
			in:             httptest.NewRequest("GET", "/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"msg":"hello world"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			h := New(nil)
			h.Root(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Fatalf("expected: %d\ngot: %d", test.expectedStatus, test.out.Code)
			}

			body := test.out.Body.String()
			if body != test.expectedBody {
				t.Fatalf("expected: %s\ngot: %s", test.expectedBody, body)
			}
		})
	}
}
