package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrObjectNotFound = errors.New("storage: object not found")
	ErrUpstream       = errors.New("storage: upstream error")
)

type Config struct {
	BaseURL        string
	ServiceRoleKey string
	Bucket         string
	HTTPTimeout    time.Duration
}

type Client interface {
	CreateSignedUploadURL(ctx context.Context, path string, ttl time.Duration) (uploadURL string, expiresAt time.Time, err error)
	Delete(ctx context.Context, path string) error
}

type SupabaseStorage struct {
	baseURL        string
	serviceRoleKey string
	bucket         string
	httpClient     *http.Client
}

func NewSupabaseStorage(cfg Config) *SupabaseStorage {
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 15 * time.Second
	}
	return &SupabaseStorage{
		baseURL:        strings.TrimRight(cfg.BaseURL, "/"),
		serviceRoleKey: cfg.ServiceRoleKey,
		bucket:         cfg.Bucket,
		httpClient:     &http.Client{Timeout: cfg.HTTPTimeout},
	}
}

type createSignedUploadResponse struct {
	URL     string `json:"url"`
	Token   string `json:"token"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (s *SupabaseStorage) CreateSignedUploadURL(ctx context.Context, path string, ttl time.Duration) (string, time.Time, error) {
	if s.baseURL == "" || s.serviceRoleKey == "" || s.bucket == "" {
		return "", time.Time{}, fmt.Errorf("storage: config incompleta (base=%q bucket=%q key_set=%v)",
			s.baseURL, s.bucket, s.serviceRoleKey != "")
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/upload/sign/%s/%s",
		s.baseURL, url.PathEscape(s.bucket), path)

	// Supabase Storage requiere body JSON, aunque sea "{}".
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte("{}")))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("storage: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", time.Time{}, fmt.Errorf("%w: status %d body=%s", ErrUpstream, resp.StatusCode, string(body))
	}

	var out createSignedUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", time.Time{}, fmt.Errorf("storage: decode: %w", err)
	}

	// Supabase devuelve `url` empezando en "/object/upload/sign/..."
	// SIN el prefijo "/storage/v1". Ese prefijo es parte de la URL pública
	// que el cliente debe hacer PUT. Lo reconstruimos.
	fullURL := out.URL
	if strings.HasPrefix(fullURL, "/") {
		fullURL = s.baseURL + "/storage/v1" + fullURL
	}
	if fullURL == "" {
		return "", time.Time{}, fmt.Errorf("%w: respuesta sin URL (token=%q path=%q)",
			ErrUpstream, out.Token, out.Path)
	}

	return fullURL, time.Now().UTC().Add(ttl), nil
}

func (s *SupabaseStorage) Delete(ctx context.Context, path string) error {
	if s.baseURL == "" || s.serviceRoleKey == "" || s.bucket == "" {
		return fmt.Errorf("storage: config incompleta")
	}

	payload, _ := json.Marshal(map[string][]string{"prefixes": {path}})
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s", s.baseURL, url.PathEscape(s.bucket))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("storage: delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("%w: delete status %d body=%s", ErrUpstream, resp.StatusCode, string(body))
	}
	return nil
}