package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DownloadToTemporaryFileFunc func(ctx context.Context, url, pattern string) (*os.File, error)

var DownloadToTemporaryFile DownloadToTemporaryFileFunc = defaultDownloadToTemporaryFile

func defaultDownloadToTemporaryFile(ctx context.Context, url, pattern string) (*os.File, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}

	if _, err := io.Copy(result, resp.Body); err != nil {
		_ = result.Close()
		_ = os.Remove(result.Name())
		return nil, fmt.Errorf("error writing to temp file: %w", err)
	}

	if _, err := result.Seek(0, io.SeekStart); err != nil {
		_ = result.Close()
		_ = os.Remove(result.Name())
		return nil, fmt.Errorf("error seeking to start of temp file: %w", err)
	}

	return result, nil
}
