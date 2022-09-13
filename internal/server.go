package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/2beens/tiny-api/pkg"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	instanceName string
	reqHandler   *RequestHandler
}

func NewServer(instanceName string) *Server {
	return &Server{
		instanceName: instanceName,
		reqHandler:   NewRequestHandler(instanceName),
	}
}

// Serve will make a server start listening on provided host and port
func (s *Server) Serve(host string, port int) {
	router := s.routerSetup()

	ipAndPort := fmt.Sprintf("%s:%d", host, port)
	httpServer := &http.Server{
		Handler:      router,
		Addr:         ipAndPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Debugf(" > server [%s] listening on: [%s]", s.instanceName, ipAndPort)
	log.Fatalf("%s: %s", s.instanceName, httpServer.ListenAndServe())
}

func (s *Server) routerSetup() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", s.reqHandler.HandlePing).Methods("GET")
	router.HandleFunc("/", s.reqHandler.HandleRootRequest).Methods("GET")

	// add a small middleware function to log request details
	router.Use(pkg.LogRequest())

	return router
}