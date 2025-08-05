package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// Captures the gin.Response and logs it before sending it
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
