package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"ph.certs.com/clm_main/auth"
	"time"
)

var lumberjackLogger = &lumberjack.Logger{
	Filename:   "middleware/log.log", //filename
	MaxSize:    100,                  // file size in MB before rotation
	MaxBackups: 10,                   // Max number of files kept before being overwritten
	MaxAge:     28,                   // Max number of days to keep the files
	Compress:   true,                 // Whether to compress log files using gzip
}
var logger = zerolog.New(lumberjackLogger).With().Timestamp().Logger()

// LoggingMiddleware is an HTTP middleware that logs request and response data, including timing and status codes.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		rec := httptest.NewRecorder()
		ctx := r.Context()
		path := r.URL.EscapedPath()
		reqData, _ := httputil.DumpRequest(r, true)

		// log time, path, request_data,
		logger := logger.Log().Timestamp().Str("path", path).Bytes("request_data", reqData)
		defer func(begin time.Time) {
			status := ww.Status()
			tookMs := time.Since(begin).Milliseconds()

			fmt.Printf("%s:%s %d\n", r.Method, path, status)
			// at the end of the call, log took, status_code, and Msg
			logger.Int64("took", tookMs).Int("status_code", status).Msgf("[%d] %s http request for %s took %dms", status, r.Method, path, tookMs)
		}(time.Now())

		// Replace "logger" with a custom type, like ContextKey("logger")
		ctx = context.WithValue(ctx, "logger", logger)
		next.ServeHTTP(rec, r.WithContext(ctx)) // this copies the recorded response to the response writer

		for k, v := range rec.Header() {
			ww.Header()[k] = v
		}
		ww.WriteHeader(rec.Code)
		_, _ = rec.Body.WriteTo(ww)

	})
}

// verifyToken validates a JWT token string and returns an error if it's invalid or cannot be parsed.
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.SecretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// CorsMiddleware is a middleware to handle Preflight calls which usually are related to CORS problems.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Currently, this only allows origins of http://localhost:4200
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
