package logging

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
)

func TestGrpcInterceptor(t *testing.T) {
	tests := []struct {
		name    string
		handler func(ctx context.Context, req interface{}) (interface{}, error)
		wantErr bool
	}{
		{
			"successful request",
			func(ctx context.Context, req interface{}) (interface{}, error) {
				return "test response", nil
			},
			false,
		},
		{
			"failed request",
			func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, assert.AnError
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}
			core, obs := observer.New(zap.InfoLevel) // Create a logger with an observer for testing
			logger := zap.New(core).Sugar()
			interceptor := GrpcInterceptor(logger)
			resp, err := interceptor(context.Background(), nil, info, tt.handler)

			// Check the response
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			// Check logged messages
			logs := obs.All()
			require.Len(t, logs, 2)
			assert.Equal(t, "Received gRPC request", logs[0].Message)
			if tt.wantErr {
				assert.Equal(t, "gRPC request failed", logs[1].Message)
			} else {
				assert.Equal(t, "gRPC request completed with response", logs[1].Message)
			}
		})
	}
}

func TestHttpMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			"successful request",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			false,
		},
		{
			"failed request",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, obs := observer.New(zap.InfoLevel) // Create a logger with an observer for testing
			logger := zap.New(core).Sugar()
			middleware := HttpMiddleware(logger)

			handler := middleware(tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Equal(t, "Internal Server Error", rec.Body.String())
			} else {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "OK", rec.Body.String())
			}

			logs := obs.All()
			require.Len(t, logs, 1)
			assert.Equal(t, "HTTP request", logs[0].Message)

			fields := logs[0].ContextMap()
			assert.Equal(t, "GET", fields["method"])
			assert.Equal(t, "/test", fields["path"])
			if tt.wantErr {
				assert.Equal(t, int64(http.StatusInternalServerError), fields["status"])
			} else {
				assert.Equal(t, int64(http.StatusOK), fields["status"])
			}
			assert.Contains(t, fields, "duration")
			assert.Contains(t, fields, "remote_addr")
		})
	}
}
