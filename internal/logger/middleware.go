package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type ctxKeyLogger struct{}

var loggerKey = ctxKeyLogger{}

func FromContext(ctx context.Context) *zerolog.Logger {
	l, ok := ctx.Value(loggerKey).(*zerolog.Logger)
	if !ok {
		panic("logger not found in context, did you forget to use logger middleware?")
	}

	return l
}

// For testing sake
type UUIDFunc func() string
type NowFunc func() time.Time

var MiddlewareRequestIDGenerator UUIDFunc = uuid.NewString
var MiddlewareNowGenerator NowFunc = time.Now

func NewMiddleware(baseLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := MiddlewareNowGenerator()

		reqLogger := baseLogger.With().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("requestID", MiddlewareRequestIDGenerator()).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.Request.UserAgent()).
			Logger()

		// Preparing to capture the response buffer
		respBuf := new(bytes.Buffer)
		writer := &bodyWriter{
			ResponseWriter: ctx.Writer,
			body:           respBuf,
		}
		ctx.Writer = writer

		// Inject logger
		ctxWithLogger := WithContext(ctx.Request.Context(), &reqLogger)
		ctx.Request = ctx.Request.WithContext(ctxWithLogger)

		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		builder := reqLogger.Info().
			Interface("request_headers", ctx.Request.Header)

		if len(requestBody) > 0 {
			builder = builder.Str("request_body", string(requestBody))
		} else {
			builder = builder.Interface("request_body", nil)
		}
		builder.Msg("Processing request")

		// Proceed with request
		ctx.Next()

		// Final log
		builder = reqLogger.Info().
			Int("status", ctx.Writer.Status()).
			Interface("response_headers", ctx.Writer.Header()).
			Float64("duration_ms", float64(MiddlewareNowGenerator().Sub(start).Microseconds())/1000.0)

		// NOTE  This only works because all responses are JSON...
		///  	 In a real app, we would do check dynamically based on the response headers
		if respBuf.Len() > 0 {
			builder = builder.RawJSON("response_body", respBuf.Bytes())
		} else {

			builder = builder.Interface("response_body", nil)
		}

		builder.
			Msg("Response sent")
	}
}

func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
