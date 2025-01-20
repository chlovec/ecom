package api

import (
	"database/sql"
	"ecom/service/user"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func InitServer(addr string, db *sql.DB) error {
	server := NewAPIServer(addr, db)
	return server.Run()
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Server started!")
	log.Println("Listening on: ", s.addr)
	return http.ListenAndServe(s.addr, router)
}
