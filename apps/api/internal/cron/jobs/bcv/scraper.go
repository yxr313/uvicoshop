package bcv

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker/v2"

	"uvico-api/internal/platform/config"
)

const (
	redisKeyRate      = "bcv:usd:micros"
	redisKeyUpdatedAt = "bcv:usd:updated_at"
	maxBodyBytes      = 2 << 20
	jobTimeout        = 45 * time.Second
)

type Scraper struct {
	cfg     config.BCVConfig
	redis   *redis.Client
	http    *http.Client
	breaker *gobreaker.CircuitBreaker[float64]
	logger  *slog.Logger
}

func New(cfg config.BCVConfig, rdb *redis.Client, logger *slog.Logger) *Scraper {
	if logger == nil {
		logger = slog.Default()
	}

	breaker := gobreaker.NewCircuitBreaker[float64](gobreaker.Settings{
		Name:        "bcv-scraper",
		MaxRequests: 1,
		Interval:    0,
		Timeout:     10 * time.Minute,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 5
		},
	})

	// El BCV usa un certificado autofirmado o fuera de la CA estándar.
	// Al ser un endpoint público (solo lectura de tasa), es seguro ignorar TLS.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Scraper{
		cfg:     cfg,
		redis:   rdb,
		http:    &http.Client{Timeout: cfg.HTTPTimeout, Transport: transport},
		breaker: breaker,
		logger:  logger,
	}
}

func (s *Scraper) Timeout() time.Duration { return jobTimeout }

func (s *Scraper) Run(ctx context.Context) error {
	rate, err := s.fetchRate(ctx)
	if err != nil {
		s.logger.Error("bcv fetch failed", "err", err)
		return fmt.Errorf("bcv fetch: %w", err)
	}

	micros := int64(math.Round(rate * 1e6))
	if micros <= 0 {
		return fmt.Errorf("bcv: invalid rate %f", rate)
	}

	return s.persist(ctx, micros, false)
}

func (s *Scraper) persist(ctx context.Context, micros int64, fallback bool) error {
	pipe := s.redis.TxPipeline()
	pipe.Set(ctx, redisKeyRate, micros, s.cfg.CacheTTL)
	pipe.Set(ctx, redisKeyUpdatedAt, time.Now().UTC().Format(time.RFC3339), s.cfg.CacheTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis persist: %w", err)
	}
	s.logger.Info("bcv rate updated", "micros", micros, "fallback", fallback)
	return nil
}

func (s *Scraper) fetchRate(ctx context.Context) (float64, error) {
	return s.breaker.Execute(func() (float64, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.cfg.ScraperURL, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Uvico-BCV-Worker/1.0)")

		resp, err := s.http.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("bcv: status %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(io.LimitReader(resp.Body, maxBodyBytes))
		if err != nil {
			return 0, fmt.Errorf("bcv: parse html: %w", err)
		}

		raw := strings.TrimSpace(doc.Find("#dolar strong").First().Text())
		if raw == "" {
			return 0, errors.New("bcv: rate not found in page")
		}

		normalized := strings.ReplaceAll(raw, ".", "")
		normalized = strings.ReplaceAll(normalized, ",", ".")

		rate, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return 0, fmt.Errorf("bcv: parse %q: %w", raw, err)
		}
		return rate, nil
	})
}