package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBodyWriter_Write(t *testing.T) {
	t.Run("should capture response body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = req
		buf := new(bytes.Buffer)

		bw := &bodyWriter{
			ResponseWriter: ctx.Writer,
			body:           buf,
		}

		bodyStr := `{"message":"test"}`

		ctx.Writer = bw
		ctx.Writer.Write([]byte(`{"message":"test"}`))

		assert.Equal(t, bodyStr, recorder.Body.String())
		assert.Equal(t, bodyStr, bw.body.String())
	})
}
