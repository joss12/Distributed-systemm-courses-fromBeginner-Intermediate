package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvstore-v3/internal/middleware"
)

func TestTimeoutMiddleware(t *testing.T) {
	logHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done() //simulate waiting and cancellation
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.Timeout(10*time.Millisecond, logHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
