package internal

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	httpServer *http.Server

	tseConn                  *grpc.ClientConn
	tseClient                tseProto.TinyStockExchangeClient
	tinyStockExchangeHandler *TinyStockExchangeHandler
}

func NewServer(
	instanceName string,
	host, port string,
	tseHost, tsePort string,
) (*Server, error) {
	tseAddr := fmt.Sprintf("%s:%s", tseHost, tsePort)
	tseConn, err := grpc.Dial(
		tseAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to dial tse connection: %w", err)
	}

	log.Debugf(" > server [%s] grpc addr: [%s]", instanceName, tseAddr)
	tseClient := tseProto.NewTinyStockExchangeClient(tseConn)

	s := &Server{
		instanceName:             instanceName,
		reqHandler:               NewRequestHandler(instanceName),
		tseConn:                  tseConn,
		tseClient:                tseClient,
		tinyStockExchangeHandler: NewTinyStockExchangeHandler(instanceName, tseClient),
	}

	ipAndPort := fmt.Sprintf("%s:%s", host, port)
	s.httpServer = &http.Server{
		Handler:      s.routerSetup(),
		Addr:         ipAndPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return s, nil
}

// Serve will make a server start listening on provided host and port
func (s *Server) Serve() {
	log.Debugf(" > server [%s] listening on: [%s]", s.instanceName, s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
		log.Warnf("[%s]: server closed", s.instanceName)
	} else {
		log.Fatalf("[%s] listen and serve: %s", s.instanceName, err)
	}
}

func (s *Server) routerSetup() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", s.reqHandler.HandlePing).Methods("GET")
	router.HandleFunc("/health", s.reqHandler.HandleHealth).Methods("GET")
	router.HandleFunc("/harakiri", s.reqHandler.HandleHarakiri).Methods("GET")
	router.HandleFunc("/", s.reqHandler.HandleRootRequest).Methods("GET")

	router.HandleFunc("/tse/stocks", s.tinyStockExchangeHandler.HandleListStocks).Methods("GET")
	router.HandleFunc("/tse/stocks", s.tinyStockExchangeHandler.HandleNewStock).Methods("POST")
	router.HandleFunc("/tse/stocks", s.tinyStockExchangeHandler.HandleUpdateStock).Methods("PATCH")
	router.HandleFunc("/tse/stocks", s.tinyStockExchangeHandler.HandleDeleteStock).Methods("DELETE")
	router.HandleFunc("/tse/deltas", s.tinyStockExchangeHandler.HandleListValueDeltas).Methods("GET")
	router.HandleFunc("/tse/deltas", s.tinyStockExchangeHandler.HandleNewValueDelta).Methods("POST")
	router.HandleFunc("/tse/status", func(w http.ResponseWriter, r *http.Request) {
		connState := s.tseConn.GetState()
		log.Printf("sending tse grpc conn state: %s", connState.String())
		pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
			Result:  "ok",
			Message: fmt.Sprintf("i[%s]: tse grpc conn state: %s", s.instanceName, connState.String()),
		})
	})

	// add a small middleware function to log request details
	router.Use(pkg.LogRequest())
	router.Use(pkg.Cors())

	return router
}

func (s *Server) Shutdown() {
	log.Debugln("server shutting down ...")
	s.tseConn.Close()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	s.httpServer.Shutdown(timeoutCtx)

	log.Debugln("bye, bye ...")
}
