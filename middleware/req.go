package middleware

import (
	"context"
	"github.com/aka-achu/watcher/logging"
	"github.com/aka-achu/watcher/response"
	"github.com/aka-achu/watcher/utility"
	"github.com/google/uuid"
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

func AuthLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		logging.Info.Printf(" [REQ] TraceID-%s Addr-%s URI-%s", traceID, r.RemoteAddr, r.RequestURI)
		if err := utility.VerifyToken(r.Header.Get("Authorization")); err != nil {
			logging.Warn.Printf(" [REQ] TraceID-%s Invalid access token. Error-%v", traceID, err)
			response.UnAuthorized(w,"invalid access token")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
