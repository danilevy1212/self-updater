package downloader

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_defaultDownloadTempFile(t *testing.T) {
	t.Run("should return a valid file", func(t *testing.T) {
		const payload = "hello world"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
		}))
		defer ts.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		f, err := defaultDownloadToTemporaryFile(ctx, ts.URL, "dltest.*")
		assert.NoError(t, err)
		assert.NotNil(t, f, "should return a valid file")
		defer func() {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}()

		off, err := f.Seek(0, io.SeekCurrent)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), off)

		got, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, payload, string(got), "should read the correct content from the file")
	})

	t.Run("should return error with non 200 status code", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		f, err := defaultDownloadToTemporaryFile(ctx, ts.URL, "dltest.*")
		assert.Error(t, err, "should return an error for non-200 status code")
		assert.Nil(t, f, "file should be nil on error")
	})
}
