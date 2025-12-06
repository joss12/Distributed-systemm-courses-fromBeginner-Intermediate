package response_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kvstore/internal/response"
)

func TestJSON(t *testing.T) {
	rr := &responseRecorder{header: http.Header{}}

	body := map[string]string{"msg": "ok"}
	response.JSON(rr, 200, body)

	if rr.status != 200 {
		t.Fatalf("expected status 200, got %d", rr.status)
	}

	var decoded map[string]string
	if err := json.Unmarshal(rr.body.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if decoded["msg"] != "ok" {
		t.Fatalf("expected msg=ok, got %s", decoded["msg"])
	}
}

func TestError(t *testing.T) {
	rr := &responseRecorder{header: http.Header{}}

	response.Error(rr, 404, "not found")

	if rr.status != 404 {
		t.Fatalf("expected 404, got %d", rr.status)
	}
}

type responseRecorder struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (rr *responseRecorder) Header() http.Header { return rr.header }
func (rr *responseRecorder) Write(b []byte) (int, error) {
	return rr.body.Write(b)
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.status = statusCode
}
