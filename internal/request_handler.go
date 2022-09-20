package internal

import (
	"fmt"
	"os"
	"net/http"
	"time"

	"github.com/2beens/tiny-api/pkg"

	log "github.com/sirupsen/logrus"
)

type RequestHandler struct {
	instanceName string
}

func NewRequestHandler(instanceName string) *RequestHandler {
	return &RequestHandler{
		instanceName: instanceName,
	}
}

func (h *RequestHandler) HandleRootRequest(w http.ResponseWriter, r *http.Request) {
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: fmt.Sprintf("Hi from instance: %s", h.instanceName),
	})
}

func (h *RequestHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Message: fmt.Sprintf("%s: PONG!", h.instanceName),
	})
}

func (h *RequestHandler) HandleHarakiri(w http.ResponseWriter, r *http.Request) {
	log.Println("killing myself ...")
	
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("killing myself..."))
	
	go func() {
		// wait a bit for the response to be sent
		time.Sleep(time.Second)
		os.Exit(1)		
	}()
}

func (h *RequestHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
