package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/watcher-io/watcher/logging"
	"github.com/watcher-io/watcher/response"
	"github.com/watcher-io/watcher/utility"
	"net/http"
)

func NoAuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [REQ] TraceID-%s Addr-%s URI-%s", traceID, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func AuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [REQ] TraceID-%s Addr-%s URI-%s", traceID, r.RemoteAddr, r.RequestURI)
		if err := utility.VerifyToken(r.Header.Get("Authorization")); err != nil {
			logging.Warn.Printf(" [REQ] TraceID-%s Invalid access token. Error-%v", traceID, err)
			response.UnAuthorized(w, "invalid access token")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
