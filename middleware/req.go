package middleware

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/google/uuid"
	"net/http"
)

func NoAuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [Req] ID-%s, URI-%s",traceID, r.RequestURI)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
