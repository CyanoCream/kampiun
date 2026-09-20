// Composition root: config, infrastruktur, rakit service (DDD), jalankan HTTP.
package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kampiun/internal/config"
	"kampiun/internal/httpapi"
	authjwt "kampiun/kernel/auth"
	"kampiun/kernel/notify"
	"kampiun/kernel/pgdb"
	"kampiun/kernel/realtime"
	"kampiun/kernel/security"
	"kampiun/service"
)

func main() {
	log.SetFlags(0)
	if err := run(config.Load()); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(lg)

	if cfg.JWTSecret == "" {
		return errors.New("JWT_SECRET wajib diisi (min 32 karakter)")
	}
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL wajib diisi")
	}

	pool, err := pgdb.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := migrate(ctx, pool); err != nil {
		return err
	}

	repo := pgdb.NewRepo(pool)
	hasher := security.NewBcryptHasher()
	hub := realtime.NewHub(cfg.RedisAddr, lg)
	defer hub.Close()

	ttl, err := time.ParseDuration(cfg.JWTAccessTTL)
	if err != nil {
		ttl = 72 * time.Hour
	}
	tokens, err := authjwt.NewHS256(cfg.JWTSecret, ttl)
	if err != nil {
		return err
	}

	if err := seed(ctx, repo, hasher, cfg.AdminEmail, cfg.SeedPassword); err != nil {
		return err
	}

	// services
	nop := notify.Nop{}
	authSvc := service.NewAuth(repo, hasher, tokens, nop)
	orgSvc := service.NewOrgService(repo)
	compSvc := service.NewCompService(repo)
	scoreSvc := service.NewScoreService(repo, hub, lg)

	// http
	authAPI := httpapi.NewAuthAPI(authSvc)
	orgAPI := httpapi.NewOrgAPI(orgSvc)
	compAPI := httpapi.NewCompAPI(compSvc, repo)
	scoreAPI := httpapi.NewScoreAPI(repo, scoreSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })

	// public
	mux.HandleFunc("POST /api/v1/auth/register", authAPI.HandleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", authAPI.HandleLogin)
	mux.HandleFunc("GET /api/v1/competitions/public", compAPI.HandlePubicSearch)
	mux.HandleFunc("GET /api/v1/matches/{id}", scoreAPI.HandleMatchGet)
	mux.HandleFunc("GET /api/v1/matches/{id}/score", scoreAPI.HandleScoreGet)
	mux.HandleFunc("GET /api/v1/matches/{id}/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(w, r, r.PathValue("id"))
	})

	// auth
	mux.Handle("GET /api/v1/orgs", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(orgAPI.HandleList)))
	mux.Handle("POST /api/v1/orgs", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(orgAPI.HandleCreate)))
	mux.Handle("POST /api/v1/competitions", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(compAPI.HandleCreate)))
	mux.Handle("GET /api/v1/orgs/{org_id}/competitions", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(compAPI.HandleListByOrg)))
	mux.Handle("POST /api/v1/competitions/{id}/participants", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(compAPI.HandleAddParticipants)))
	mux.Handle("POST /api/v1/matches/{id}/tap", httpapi.MiddlewareAuth(tokens, http.HandlerFunc(scoreAPI.HandleTap)))

	srv := &http.Server{Addr: cfg.Addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	lg.Info("kampiun listening", "addr", cfg.Addr)
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
	case err := <-errc:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	shCtx, c2 := context.WithTimeout(ctx, 5*time.Second)
	defer c2()
	return srv.Shutdown(shCtx)
}
