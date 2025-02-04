package http

import (
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type middleware func(logger *slog.Logger, handler http.Handler) http.Handler

type Server struct {
	tenantHandler TenantHandler
	router        http.Handler
	logger        *slog.Logger
	middlewares   []middleware
}

func NewServer(pool *pgxpool.Pool, logger *slog.Logger) *Server {
	tenantStore := postgres.NewTenantStore(pool)
	userStore := postgres.NewUserStore(pool)

	tenantHandler := NewTenantHandler(logger, tenantStore, userStore)

	router := http.NewServeMux()

	server := &Server{
		tenantHandler: *tenantHandler,
		router:        router,
		logger:        logger,
	}

	server.registerRoutes(router)
	return server
}

func (s *Server) registerRoutes(router *http.ServeMux) {
	router.Handle("/api/tenants/", s.tenantHandler)
}

func (s *Server) Use(m middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) registerGlobalMiddlewares() http.Handler {
	handler := s.router
	for _, m := range s.middlewares {
		handler = m(s.logger, handler)
	}
	return handler
}

func (s *Server) ListenAndServe(port string) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.registerGlobalMiddlewares(),
	}
	s.logger.Info("server is running", slog.String("port", port))
	return server.ListenAndServe()
}
