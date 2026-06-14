package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func NormalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	if !strings.EqualFold(parsed.Host, "md.archlinux.org") {
		return rawURL
	}
	if strings.HasSuffix(parsed.Path, "/download") || strings.HasSuffix(parsed.Path, "/raw") {
		return rawURL
	}

	cleanPath := strings.TrimSuffix(parsed.Path, "/")
	if cleanPath == "" || cleanPath == "/" {
		return rawURL
	}
	parsed.Path = path.Clean(cleanPath) + "/download"
	return parsed.String()
}

func FetchURL(ctx context.Context, rawURL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "aur-scanner/0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
