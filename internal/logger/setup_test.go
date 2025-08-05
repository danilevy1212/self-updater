package logger

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {
	// Gin test mode
	gin.SetMode(gin.TestMode)

	// Freeze time
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2025, 3, 27, 12, 0, 0, 0, time.UTC)
	}
	defer func() {
		zerolog.TimestampFunc = time.Now
	}()
	nowFunc := func() time.Time {
		return time.Date(2025, 3, 27, 12, 0, 0, 0, time.UTC)
	}
	MiddlewareNowGenerator = nowFunc

	// Freeze UUID generation
	uuidGen := func() string {
		return "6460a629-fab4-4eaf-8de1-bfb4f39ddbbc"
	}
	MiddlewareRequestIDGenerator = uuidGen

	v := m.Run()

	snaps.Clean(m)
	os.Exit(v)
}
