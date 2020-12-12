package middleware

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/utility"
	"github.com/google/uuid"
	"net/http"
)

// NoAuthLogging is a middleware which will omit the request auth validation
// and generate a request tracing id.
func NoAuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generating a request tracing id
		traceID := uuid.New().String()
		// Storing the generated traceID in the context store against the key "trace_id"
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [Req] ID-%s, URI-%s", traceID, r.RequestURI)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthLogging is a middleware which will validate the request based on the
// Authorization token present in the request header and generate a request tracing id.
func AuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generating a request tracing id
		traceID := uuid.New().String()
		// Storing the generated traceID in the context store against the key "trace_id"
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [Req] ID-%s, URI-%s", traceID, r.RequestURI)
		// Validating the Authorization token present in the request header
		if err := utility.VerifyToken(r.Header.Get("Authorization")); err != nil {
			logging.Warn.Printf(" [Req] Invalid access token. Error-%b TraceID-%s", err, traceID)
			response.UnAuthorized(w,"Invalid access token")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
