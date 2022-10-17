package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/2beens/tiny-api/pkg"
	tseProto "github.com/2beens/tiny-stock-exchange-proto"

	log "github.com/sirupsen/logrus"
)

type TinyStockExchangeHandler struct {
	instanceName string
	tseClient    tseProto.TinyStockExchangeClient
}

func NewTinyStockExchangeHandler(instanceName string) *TinyStockExchangeHandler {
	return &TinyStockExchangeHandler{
		instanceName: instanceName,
	}
}

func (h *TinyStockExchangeHandler) HandleNewStock(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Errorf("add new stock failed, parse form error: %s", err)
		http.Error(w, "parse form error", http.StatusInternalServerError)
		return
	}

	ticker := r.Form.Get("ticker")
	if ticker == "" {
		http.Error(w, "error, ticker empty", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "error, name empty", http.StatusBadRequest)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	h.tseClient.NewStock(timeoutCtx, &tseProto.Stock{
		Ticker: ticker,
		Name:   name,
	})

	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: fmt.Sprintf("i[%s]: new stock %s added", h.instanceName, ticker),
	})
}

func (h *TinyStockExchangeHandler) HandleNewValueDelta(w http.ResponseWriter, r *http.Request) {
	// TODO:

	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: fmt.Sprintf("Hi from instance: %s", h.instanceName),
	})
}
