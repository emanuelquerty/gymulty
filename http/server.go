package http

import (
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/multency/postgres"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	tenantHandler TenantHandler
	router        http.Handler
	logger        *slog.Logger
}

func NewServer(conn *pgx.Conn, logger *slog.Logger) *Server {
	tenantStore := postgres.NewTenantStore(conn)
	userStore := postgres.NewUserStore(conn)

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

func (s *Server) ListenAndServe(port string) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}
	s.logger.Info("server is running", slog.String("port", port))
	return server.ListenAndServe()
}
