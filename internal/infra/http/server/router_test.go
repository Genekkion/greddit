package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"greddit/internal/test"
)

func TestHttpServer_Health(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	health(w, r)

	res := w.Result()
	defer res.Body.Close()

	test.AssertEqual(t, "Unexpected status code", http.StatusOK, res.StatusCode)
	var got struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(res.Body).Decode(&got)
	test.NilErr(t, err)
	test.AssertEqual(t, "Unexpected status", "healthy", got.Status)
}
