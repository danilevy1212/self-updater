package logger

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_middleware_NewMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	logger := zerolog.
		New(&logBuf).
		With().
		Timestamp().
		Logger()

	router := gin.Default()
	router.Use(NewMiddleware(&logger))
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	t.Run("should add logs add the beginning and end of the request lifetime", func(t *testing.T) {
		responseReader := strings.NewReader(`Hey`)
		req := httptest.NewRequest(http.MethodGet, "/", responseReader)
		req.Header.Set("User-Agent", "TestAgent/1.0")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "code should be 200")
		snaps.MatchSnapshot(t, logBuf.String())
	})

	t.Run("should omit the request and response body when not provided", func(t *testing.T) {
		router.DELETE("/", func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		})
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("User-Agent", "TestAgent/1.0")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code, "code should be 204")
		snaps.MatchSnapshot(t, logBuf.String())
	})
}

func Test_FromContext(t *testing.T) {
	tests := []struct {
		name  string
		setup bool
	}{
		{"should return a logger when set", true},
		{"should panic when not set", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup {
				logger := zerolog.New(zerolog.NewConsoleWriter())
				ctx := WithContext(context.Background(), &logger)
				assert.NotPanics(t, func() {
					newLogger := FromContext(ctx)
					assert.NotNil(t, newLogger, "should not be nil")
					assert.Same(t, &logger, newLogger, "loggers are not the same")
				})
			} else {
				assert.Panics(t, func() {
					FromContext(context.Background())
				})
			}
		})
	}
}

func Test_WithContext(t *testing.T) {
	t.Run("should return a context with a set logger", func(t *testing.T) {
		logger := zerolog.New(zerolog.NewConsoleWriter())
		context := WithContext(context.Background(), &logger)

		newLogger, ok := context.Value(loggerKey).(*zerolog.Logger)

		assert.True(t, ok, "should not be nil")
		assert.Same(t, &logger, newLogger, "loggers are not the same")
	})
}
