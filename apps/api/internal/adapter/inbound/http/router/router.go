package router

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	bcvhttp "uvico-api/internal/adapter/inbound/http/bcv"
	billinghttp "uvico-api/internal/adapter/inbound/http/billing"
	carthttp "uvico-api/internal/adapter/inbound/http/cart"
	cataloghttp "uvico-api/internal/adapter/inbound/http/catalog"
	imagehttp "uvico-api/internal/adapter/inbound/http/image"
	"uvico-api/internal/adapter/inbound/http/middleware"
	orderhttp "uvico-api/internal/adapter/inbound/http/order"
	paymentmethodhttp "uvico-api/internal/adapter/inbound/http/paymentmethod"
	storefronthttp "uvico-api/internal/adapter/inbound/http/storefront"
	tenanthttp "uvico-api/internal/adapter/inbound/http/tenant"
	webhookhttp "uvico-api/internal/adapter/inbound/http/webhook"
	"uvico-api/internal/adapter/outbound/postgres"
	redisadapter "uvico-api/internal/adapter/outbound/redis"
	"uvico-api/internal/adapter/outbound/storage"
	appbilling "uvico-api/internal/application/billing"
	appcart "uvico-api/internal/application/cart"
	appcatalog "uvico-api/internal/application/catalog"
	appimage "uvico-api/internal/application/image"
	apporder "uvico-api/internal/application/order"
	apppayment "uvico-api/internal/application/payment"
	apppaymentmethod "uvico-api/internal/application/paymentmethod"
	appstorefront "uvico-api/internal/application/storefront"
	apptenant "uvico-api/internal/application/tenant"
	billdomain "uvico-api/internal/domain/billing"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/platform/config"
	"uvico-api/internal/platform/database"
)

type Dependencies struct {
	Pool   *database.Pool
	Redis  *redis.Client
	Logger *slog.Logger
	Cfg    *config.Config
	Plans  map[billdomain.PlanID]billdomain.Plan
}

func New(deps *Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer(deps.Logger))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Logger(deps.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.Cfg.Security.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Tenant-ID", "X-Request-ID", "Idempotency-Key", "X-Session-ID", "X-Order-Token"},
		ExposedHeaders:   []string{"X-Request-ID", "Idempotency-Key", "Idempotent-Replay", "X-Session-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		if err := deps.Pool.Ping(ctx); err != nil {
			deps.Logger.Warn("readiness: database unreachable", "err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"unavailable","dependency":"postgres"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})
	r.Handle("/metrics", promhttp.Handler())

	// ─── Adaptadores ───────────────────────────────────────────────
	storefrontRepo := postgres.NewStorefrontRepo(deps.Pool)
	bcvRepo := redisadapter.NewBCVRepo(deps.Redis)
	storefrontSvc := appstorefront.NewService(storefrontRepo, bcvRepo)
	storefrontResolver := storefronthttp.NewSlugResolver(storefrontSvc)
	bcvH := bcvhttp.New(storefrontSvc, deps.Logger)

	tenantRepo := postgres.NewTenantRepo(deps.Pool)
	subRepo := postgres.NewSubscriptionRepo(deps.Pool)
	billingSvc := appbilling.NewService(tenantRepo, subRepo, deps.Plans, deps.Logger)

	catalogRepo := postgres.NewCatalogRepo(deps.Pool)
	catalogSvc := appcatalog.NewService(catalogRepo)

	invRepo := postgres.NewInventoryRepo(deps.Pool)
	orderRepo := postgres.NewOrderRepo(deps.Pool)
	payRepo := postgres.NewPaymentRepo()
	customerRepo := postgres.NewCustomerRepo(deps.Pool)
	attributionRepo := postgres.NewAttributionRepo()
	pmRepo := postgres.NewPaymentMethodRepo(deps.Pool)
	idemStore := redisadapter.NewIdempotencyStore(deps.Redis)

	cartRepo := redisadapter.NewCartRepo(deps.Redis)
	cartSvc := appcart.NewService(cartRepo, catalogRepo)

	storageClient := storage.NewSupabaseStorage(storage.Config{
		BaseURL:        deps.Cfg.Supabase.URL,
		ServiceRoleKey: deps.Cfg.Supabase.ServiceRoleKey,
		Bucket:         deps.Cfg.Storage.Bucket,
		HTTPTimeout:    15 * time.Second,
	})
	imageRepo := postgres.NewImageRepo(deps.Pool)
	imageSvc := appimage.NewService(imageRepo, catalogRepo, storageClient, appimage.Options{
		MaxFileSizeBytes: deps.Cfg.Storage.MaxFileSizeBytes,
		AllowedMIMETypes: deps.Cfg.Storage.AllowedMIMETypes,
		PresignedURLTTL:  deps.Cfg.Storage.PresignedURLTTL,
	})
	imageH := imagehttp.New(imageSvc, deps.Logger)

	orderSvc := apporder.NewService(
		deps.Pool,
		orderRepo,
		invRepo,
		customerRepo,
		attributionRepo,
		pmRepo,
		idemStore,
		deps.Logger,
		apporder.ServiceOptions{
			Split: orderdomain.SplitConfig{
				UvicoCommissionBps: deps.Cfg.Custody.UvicoCommissionBPS,
				ProviderCostBps:    deps.Cfg.Custody.ProviderCostDefaultBPS,
				GatewayFeeBps:      deps.Cfg.Custody.GatewayFeeDefaultBPS,
			},
			CustodyEnabled: deps.Cfg.Custody.CustodyEnabled,
		},
	)

	orderH := orderhttp.New(orderSvc, deps.Logger)
	cartH := carthttp.New(cartSvc, deps.Logger)

	pmSvc := apppaymentmethod.NewService(pmRepo)
	pmH := paymentmethodhttp.New(pmSvc, deps.Logger)

	storefrontH := storefronthttp.New(
		storefrontSvc,
		orderSvc,
		cartSvc,
		pmSvc,
		deps.Logger,
		[]byte(deps.Cfg.Security.HMACSecret),
	)

	webhookSvc := apppayment.NewWebhookService(deps.Pool, orderRepo, invRepo, payRepo, deps.Logger)

	tenantSvc := apptenant.NewService(deps.Pool, deps.Logger)
	jwtSvc := middleware.NewJWTService(deps.Cfg.Supabase, deps.Pool, deps.Redis, deps.Logger)

	// ─── Rutas ─────────────────────────────────────────────────────
	r.Route(deps.Cfg.Server.BasePath, func(r chi.Router) {
		// Público (rate limit por IP)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(deps.Redis, 60, time.Minute))
			r.Mount("/billing", billinghttp.New(billingSvc, deps.Cfg, deps.Logger).Routes())
			r.Mount("/storefront", storefrontH.Routes(storefrontResolver))
			r.Get("/bcv/rate", bcvH.GetRate)
		})

		// Webhooks (HMAC, sin rate limit)
		r.Mount("/webhooks", webhookhttp.New(webhookSvc, deps.Cfg, deps.Logger).Routes())

		// Autenticadas (JWT)
		r.Group(func(r chi.Router) {
			r.Use(jwtSvc.Middleware())
			r.Mount("/tenants", tenanthttp.New(tenantSvc, jwtSvc, deps.Logger).Routes())

			r.Group(func(r chi.Router) {
				r.Use(middleware.RateLimit(deps.Redis, 100, time.Second))
				r.Mount("/products", cataloghttp.New(catalogSvc, deps.Logger).Routes())
				r.Mount("/orders", orderH.Routes())
				r.Mount("/cart", cartH.Routes())
				r.Mount("/uploads", imageH.UploadRoutes())
				r.Mount("/store/payment-methods", pmH.MerchantRoutes())

				r.Route("/products/{productID}/images", func(r chi.Router) {
					r.Mount("/", imageH.ImageRoutes())
				})
			})
		})
	})

	return r
}