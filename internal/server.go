package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/2beens/tiny-api/pkg"
	tseProto "github.com/2beens/tiny-stock-exchange-proto"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	instanceName string
	reqHandler   *RequestHandler

	tseClient                tseProto.TinyStockExchangeClient
	tinyStockExchangeHandler *TinyStockExchangeHandler
}

func NewServer(instanceName string, tseClient tseProto.TinyStockExchangeClient) *Server {
	return &Server{
		instanceName:             instanceName,
		reqHandler:               NewRequestHandler(instanceName),
		tinyStockExchangeHandler: NewTinyStockExchangeHandler(instanceName, tseClient),
	}
}

// Serve will make a server start listening on provided host and port
func (s *Server) Serve(host, port string) {
	router := s.routerSetup()

	ipAndPort := fmt.Sprintf("%s:%s", host, port)
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
	router.HandleFunc("/health", s.reqHandler.HandleHealth).Methods("GET")
	router.HandleFunc("/harakiri", s.reqHandler.HandleHarakiri).Methods("GET")
	router.HandleFunc("/", s.reqHandler.HandleRootRequest).Methods("GET")

	router.HandleFunc("/tse/stocks", s.tinyStockExchangeHandler.HandleNewStock).Methods("POST")
	router.HandleFunc("/tse/deltas", s.tinyStockExchangeHandler.HandleNewValueDelta).Methods("POST")

	// add a small middleware function to log request details
	router.Use(pkg.LogRequest())
	router.Use(pkg.Cors())

	return router
}
