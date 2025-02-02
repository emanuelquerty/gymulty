package http

import (
	"fmt"
	"net/http"

	"github.com/emanuelquerty/multency/postgres"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	tenantHandler TenantHandler
	router        http.Handler
}

func NewServer(conn *pgx.Conn) *Server {
	tenantStore := postgres.NewTenantStore(conn)
	userStore := postgres.NewUserStore(conn)

	tenantHandler := NewTenantHandler(tenantStore, userStore)

	router := http.NewServeMux()

	server := &Server{
		tenantHandler: *tenantHandler,
		router:        router,
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

	fmt.Println("Server listening on port", port)
	return server.ListenAndServe()
}
