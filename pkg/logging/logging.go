package logging

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var loggerSet *zap.Logger

func Initialize(debug bool) error {
	var err error
	var logger *zap.Logger
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	loggerSet = logger
	defer loggerSet.Sync()
	return nil
}

func NewLogger(name string) *zap.SugaredLogger {
	return loggerSet.Named(name).Sugar()
}

func GrpcInterceptor(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infow("Received gRPC request", "method", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Errorw("gRPC request failed", "method", info.FullMethod, "error", err)
		} else {
			logger.Infow("gRPC request completed with response", "method", info.FullMethod, "response", resp)
		}
		return resp, err
	}
}

func HttpMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrappedWriter := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(wrappedWriter, r) // Call the next handler
			logger.Infow("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"remote_addr", r.RemoteAddr,
				"status", wrappedWriter.status,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
