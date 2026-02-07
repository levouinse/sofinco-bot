package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type YouTubeDownloader struct {
	apiKey  string
	baseURL string
}

type YouTubeResult struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Title       string `json:"title"`
		Duration    string `json:"duration"`
		Views       string `json:"views"`
		Author      string `json:"author"`
		Thumbnail   string `json:"thumbnail"`
		MP3         string `json:"mp3"`
		MP4         string `json:"mp4"`
		Description string `json:"description"`
	} `json:"result"`
}

func NewYouTubeDownloader(apiKey, baseURL string) *YouTubeDownloader {
	return &YouTubeDownloader{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (yt *YouTubeDownloader) Download(query string) (*YouTubeResult, error) {
	endpoint := fmt.Sprintf("%s/api/download/yt?url=%s&apikey=%s",
		yt.baseURL,
		url.QueryEscape(query),
		yt.apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result YouTubeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result, nil
}

func (yt *YouTubeDownloader) Search(query string) (*YouTubeResult, error) {
	endpoint := fmt.Sprintf("%s/api/search/yts?query=%s&apikey=%s",
		yt.baseURL,
		url.QueryEscape(query),
		yt.apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result YouTubeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
