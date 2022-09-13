package internal

import (
	"fmt"
	"net/http"

	"github.com/2beens/tiny-api/pkg"
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
