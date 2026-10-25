// Package config carga y valida la configuración desde variables de entorno.
// Fail-fast: errores de configuración impiden el arranque.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

type Config struct {
	Env           Environment
	Service       ServiceConfig
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Supabase      SupabaseConfig
	Auth          AuthConfig
	Observability ObservabilityConfig
	Security      SecurityConfig
	BCV           BCVConfig
	Plans         PlansConfig
	Custody       CustodyConfig
	Tenancy       TenancyConfig
	Messaging     MessagingConfig
	Webhooks      WebhooksConfig
	Storage       StorageConfig
}

type ServiceConfig struct {
	Name    string
	Version string
	Region  string
}

type ServerConfig struct {
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
	BasePath          string
}

type DatabaseConfig struct {
	URL              string
	MigrationsURL    string
	MaxConns         int32
	MinConns         int32
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
	HealthCheck      time.Duration
	ConnectTimeout   time.Duration
	StatementTimeout time.Duration
	SSLMode          string
}

type RedisConfig struct {
	URL          string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type SupabaseConfig struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	JWTSecret      string
}

type AuthConfig struct {
	JWTIssuer       string
	JWTAudience     string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type ObservabilityConfig struct {
	LogLevel        string
	OTelEndpoint    string
	OTelServiceName string
	OTelSampleRate  float64
	SentryDSN       string
	MetricsPort     string
}

type SecurityConfig struct {
	CORSAllowedOrigins []string
	RateLimitRedisKey  string
	EnforceHTTPS       bool
	HMACSecret         string
}

type BCVConfig struct {
	ScraperURL  string
	ScraperCron string
	HTTPTimeout time.Duration
	CacheTTL    time.Duration
}

type PlansConfig struct {
	FounderMonthlyCents      int64
	FounderTrialCents        int64
	EssentialMonthlyCents    int64
	ProfessionalMonthlyCents int64
	EnterpriseMonthlyCents   int64
	AnnualDiscountBPS        int64
}

// CustodyConfig — una sola declaración. No duplicar.
type CustodyConfig struct {
	UvicoCommissionBPS     int64
	GatewayFeeDefaultBPS   int64
	ProviderCostDefaultBPS int64
	PayoutWeeklyDOW        int
	PayoutFortnightlyDOW   int
	PayoutMinCents         int64
	CustodyEnabled         bool
}

type TenancyConfig struct {
	DefaultTenantID  string
	EnforceIsolation bool
	HeaderName       string
}

type MessagingConfig struct {
	Brokers            []string
	OrderEventsTopic   string
	PaymentEventsTopic string
	ConsumerGroup      string
	OutboxBatchSize    int
	OutboxPollInterval time.Duration
}

type PaymentGatewayConfig struct {
	Name          string
	BaseURL       string
	APIKey        string
	WebhookSecret string
	WebhookHeader string
}

type WebhooksConfig struct {
	Wayu    PaymentGatewayConfig
	Veripay PaymentGatewayConfig
	VEXPay  PaymentGatewayConfig
	Payku   PaymentGatewayConfig
	Pabilo  PaymentGatewayConfig
}

type StorageConfig struct {
	Bucket           string
	MaxFileSizeBytes int64
	AllowedMIMETypes []string
	PresignedURLTTL  time.Duration
}

func Load() (*Config, error) {
	env := Environment(strings.ToLower(getEnv("APP_ENV", string(EnvDevelopment))))
	if env != EnvProduction {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Env: env,
		Service: ServiceConfig{
			Name:    getEnv("SERVICE_NAME", "uvico-api"),
			Version: getEnv("SERVICE_VERSION", "dev"),
			Region:  getEnv("SERVICE_REGION", "local"),
		},
		Server: ServerConfig{
			Port:              getEnv("PORT", "8080"),
			ReadTimeout:       getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			ReadHeaderTimeout: getEnvDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
			WriteTimeout:      getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:       getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout:   getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 25*time.Second),
			MaxHeaderBytes:    getEnvInt("SERVER_MAX_HEADER_BYTES", 1<<20),
			BasePath:          getEnv("SERVER_BASE_PATH", "/v1"),
		},
		Database: DatabaseConfig{
			URL:              getEnv("DATABASE_URL", ""),
			MigrationsURL:    getEnv("MIGRATIONS_URL", getEnv("DATABASE_URL", "")),
			MaxConns:         int32(getEnvInt("DB_MAX_CONNS", 25)),
			MinConns:         int32(getEnvInt("DB_MIN_CONNS", 5)),
			MaxConnLifetime:  getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			MaxConnIdleTime:  getEnvDuration("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			HealthCheck:      getEnvDuration("DB_HEALTH_CHECK", 1*time.Minute),
			ConnectTimeout:   getEnvDuration("DB_CONNECT_TIMEOUT", 5*time.Second),
			StatementTimeout: getEnvDuration("DB_STATEMENT_TIMEOUT", 5*time.Second),
			SSLMode:          getEnv("DB_SSL_MODE", "require"),
		},
		Redis: RedisConfig{
			URL:          getEnv("REDIS_URL", ""),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvInt("REDIS_DB", 0),
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 50),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 10),
			DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 3*time.Second),
			ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		Supabase: SupabaseConfig{
			URL:            getEnv("SUPABASE_URL", ""),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			JWTSecret:      getEnv("SUPABASE_JWT_SECRET", ""),
		},
		Auth: AuthConfig{
			JWTIssuer:       getEnv("AUTH_JWT_ISSUER", "uvico"),
			JWTAudience:     getEnv("AUTH_JWT_AUDIENCE", "uvico-api"),
			AccessTokenTTL:  getEnvDuration("AUTH_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("AUTH_REFRESH_TTL", 7*24*time.Hour),
		},
		Observability: ObservabilityConfig{
			LogLevel:        getEnv("LOG_LEVEL", "info"),
			OTelEndpoint:    getEnv("OTEL_ENDPOINT", ""),
			OTelServiceName: getEnv("OTEL_SERVICE_NAME", "uvico-api"),
			OTelSampleRate:  getEnvFloat("OTEL_SAMPLE_RATE", 0.1),
			SentryDSN:       getEnv("SENTRY_DSN", ""),
			MetricsPort:     getEnv("METRICS_PORT", "9090"),
		},
		Security: SecurityConfig{
			CORSAllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			RateLimitRedisKey:  getEnv("RATE_LIMIT_REDIS_KEY", "uvico:rl"),
			EnforceHTTPS:       getEnvBool("ENFORCE_HTTPS", env == EnvProduction),
			HMACSecret:         getEnv("HMAC_SECRET", ""),
		},
		BCV: BCVConfig{
			ScraperURL:  getEnv("BCV_SCRAPER_URL", "https://www.bcv.org.ve/"),
			ScraperCron: getEnv("BCV_SCRAPER_CRON", "0 * * * *"),
			HTTPTimeout: getEnvDuration("BCV_HTTP_TIMEOUT", 10*time.Second),
			CacheTTL:    getEnvDuration("BCV_CACHE_TTL", 1*time.Hour),
		},
		Plans: PlansConfig{
			FounderMonthlyCents:      getEnvInt64("PLAN_FOUNDER_MONTHLY_CENTS", 900),
			FounderTrialCents:        getEnvInt64("PLAN_FOUNDER_TRIAL_CENTS", 100),
			EssentialMonthlyCents:    getEnvInt64("PLAN_ESSENTIAL_MONTHLY_CENTS", 1999),
			ProfessionalMonthlyCents: getEnvInt64("PLAN_PROFESSIONAL_MONTHLY_CENTS", 4999),
			EnterpriseMonthlyCents:   getEnvInt64("PLAN_ENTERPRISE_MONTHLY_CENTS", 9999),
			AnnualDiscountBPS:        getEnvInt64("PLAN_ANNUAL_DISCOUNT_BPS", 8300),
		},
		Custody: CustodyConfig{
			UvicoCommissionBPS:     getEnvInt64("UVICO_COMMISSION_BPS", 500),
			GatewayFeeDefaultBPS:   getEnvInt64("GATEWAY_FEE_DEFAULT_BPS", 350),
			ProviderCostDefaultBPS: getEnvInt64("PROVIDER_COST_DEFAULT_BPS", 5000),
			PayoutWeeklyDOW:        getEnvInt("PAYOUT_WEEKLY_DOW", 5),
			PayoutFortnightlyDOW:   getEnvInt("PAYOUT_FORTNIGHTLY_DOW", 5),
			PayoutMinCents:         getEnvInt64("PAYOUT_MIN_CENTS", 1000),
			CustodyEnabled:         getEnvBool("CUSTODY_ENABLED", false),
		},
		Tenancy: TenancyConfig{
			DefaultTenantID:  getEnv("DEFAULT_TENANT_ID", ""),
			EnforceIsolation: getEnvBool("ENFORCE_TENANT_ISOLATION", true),
			HeaderName:       getEnv("TENANT_HEADER", "X-Tenant-ID"),
		},
		Messaging: MessagingConfig{
			Brokers:            getEnvSlice("KAFKA_BROKERS", nil),
			OrderEventsTopic:   getEnv("KAFKA_ORDER_TOPIC", "uvico.orders.v1"),
			PaymentEventsTopic: getEnv("KAFKA_PAYMENT_TOPIC", "uvico.payments.v1"),
			ConsumerGroup:      getEnv("KAFKA_CONSUMER_GROUP", "uvico-worker"),
			OutboxBatchSize:    getEnvInt("OUTBOX_BATCH_SIZE", 100),
			OutboxPollInterval: getEnvDuration("OUTBOX_POLL_INTERVAL", 2*time.Second),
		},
		Webhooks: WebhooksConfig{
			Wayu:    loadGateway("WAYU"),
			Veripay: loadGateway("VERIPAY"),
			VEXPay:  loadGateway("VEXPAY"),
			Payku:   loadGateway("PAYKU"),
			Pabilo:  loadGateway("PABILO"),
		},
		Storage: StorageConfig{
			Bucket:           getEnv("SUPABASE_STORAGE_BUCKET", "products"),
			MaxFileSizeBytes: getEnvInt64("UPLOAD_MAX_SIZE_BYTES", 5<<20),
			AllowedMIMETypes: getEnvSlice("UPLOAD_ALLOWED_MIMES", []string{
				"image/jpeg", "image/png", "image/webp", "image/avif",
			}),
			PresignedURLTTL: getEnvDuration("UPLOAD_PRESIGNED_TTL", 5*time.Minute),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

func (c *Config) IsProduction() bool { return c.Env == EnvProduction }

func (c *Config) Validate() error {
	var errs []error

	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		errs = append(errs, fmt.Errorf("PORT inválido (%q)", c.Server.Port))
	}
	if c.Database.MaxConns < c.Database.MinConns {
		errs = append(errs, errors.New("DB_MAX_CONNS < DB_MIN_CONNS"))
	}
	if c.Server.ReadHeaderTimeout <= 0 {
		errs = append(errs, errors.New("SERVER_READ_HEADER_TIMEOUT debe ser > 0"))
	}
	if c.Server.ShutdownTimeout < 10*time.Second {
		errs = append(errs, errors.New("SERVER_SHUTDOWN_TIMEOUT debe ser >= 10s"))
	}

	if c.IsProduction() {
		for name, val := range map[string]string{
			"DATABASE_URL":              c.Database.URL,
			"REDIS_URL":                 c.Redis.URL,
			"SUPABASE_URL":              c.Supabase.URL,
			"SUPABASE_SERVICE_ROLE_KEY": c.Supabase.ServiceRoleKey,
			"SUPABASE_JWT_SECRET":       c.Supabase.JWTSecret,
			"OTEL_ENDPOINT":             c.Observability.OTelEndpoint,
			"HMAC_SECRET":               c.Security.HMACSecret,
			"SUPABASE_STORAGE_BUCKET":   c.Storage.Bucket,
		} {
			if val == "" {
				errs = append(errs, fmt.Errorf("%s requerido en producción", name))
			}
		}
		if contains(c.Security.CORSAllowedOrigins, "*") {
			errs = append(errs, errors.New("CORS_ALLOWED_ORIGINS no puede ser '*' en producción"))
		}
		if c.Tenancy.DefaultTenantID != "" {
			errs = append(errs, errors.New("DEFAULT_TENANT_ID debe estar vacío en producción"))
		}
	}

	for name, raw := range map[string]string{
		"DATABASE_URL":    c.Database.URL,
		"MIGRATIONS_URL":  c.Database.MigrationsURL,
		"REDIS_URL":       c.Redis.URL,
		"SUPABASE_URL":    c.Supabase.URL,
		"BCV_SCRAPER_URL": c.BCV.ScraperURL,
	} {
		if raw == "" {
			continue
		}
		if _, err := url.ParseRequestURI(raw); err != nil {
			errs = append(errs, fmt.Errorf("%s inválida: %w", name, err))
		}
	}

	if c.Plans.AnnualDiscountBPS <= 0 || c.Plans.AnnualDiscountBPS > 10000 {
		errs = append(errs, errors.New("PLAN_ANNUAL_DISCOUNT_BPS fuera de rango"))
	}
	if c.Custody.UvicoCommissionBPS < 0 || c.Custody.UvicoCommissionBPS > 10000 {
		errs = append(errs, errors.New("UVICO_COMMISSION_BPS fuera de rango"))
	}
	if c.Storage.MaxFileSizeBytes <= 0 || c.Storage.MaxFileSizeBytes > 50<<20 {
		errs = append(errs, errors.New("UPLOAD_MAX_SIZE_BYTES fuera de rango (1..50MB)"))
	}
	if len(c.Storage.AllowedMIMETypes) == 0 {
		errs = append(errs, errors.New("UPLOAD_ALLOWED_MIMES no puede estar vacío"))
	}
	if c.Storage.PresignedURLTTL <= 0 || c.Storage.PresignedURLTTL > 60*time.Minute {
		errs = append(errs, errors.New("UPLOAD_PRESIGNED_TTL fuera de rango (1m..60m)"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return fallback
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func loadGateway(prefix string) PaymentGatewayConfig {
	return PaymentGatewayConfig{
		Name:          strings.ToLower(prefix),
		BaseURL:       getEnv(prefix+"_BASE_URL", ""),
		APIKey:        getEnv(prefix+"_API_KEY", ""),
		WebhookSecret: getEnv(prefix+"_WEBHOOK_SECRET", ""),
		WebhookHeader: getEnv(prefix+"_WEBHOOK_HEADER", "X-Signature"),
	}
}