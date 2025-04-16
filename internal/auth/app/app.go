package app

import (
	"context"
	"encoding/json"
	"github.com/Chamistery/TestTask/internal/auth/auth_http"
	"github.com/Chamistery/TestTask/internal/auth/closer"
	"github.com/Chamistery/TestTask/internal/auth/config"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Chamistery/TestTask/internal/auth/logger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPAuthServer,
	}
	logger.Init(getCore(getAtomicLevel()))
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPAuthServer(ctx context.Context) error {
	mux := http.NewServeMux()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	authHandler := a.serviceProvider.HTTPAuthImpl(ctx)

	mux.HandleFunc("/Create", func(w http.ResponseWriter, r *http.Request) {
		a.handleCreate(w, r, authHandler)
	})

	mux.HandleFunc("/Refresh", func(w http.ResponseWriter, r *http.Request) {
		a.handleRefresh(w, r, authHandler)
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (a *App) handleCreate(w http.ResponseWriter, r *http.Request, handler *auth_http.Implementation) {
	if r.Method == http.MethodPost {
		ip := getClientIP(r)
		var req model.CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := handler.Create(r.Context(), &req, ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) handleRefresh(w http.ResponseWriter, r *http.Request, handler *auth_http.Implementation) {
	if r.Method == http.MethodPatch {
		ip := getClientIP(r)
		var req model.RefreshRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := handler.Refresh(r.Context(), &req, ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.httpServer.Addr)
	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPAuthConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}
