package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/2beens/tiny-api/pkg"
	tseProto "github.com/2beens/tiny-stock-exchange-proto"

	log "github.com/sirupsen/logrus"
)

type TinyStockExchangeHandler struct {
	instanceName string
	tseClient    tseProto.TinyStockExchangeClient
}

func NewTinyStockExchangeHandler(
	instanceName string,
	tseClient tseProto.TinyStockExchangeClient,
) *TinyStockExchangeHandler {
	return &TinyStockExchangeHandler{
		instanceName: instanceName,
		tseClient:    tseClient,
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
	res, err := h.tseClient.NewStock(timeoutCtx, &tseProto.Stock{
		Ticker: ticker,
		Name:   name,
	})
	if err != nil {
		log.Errorf("add stock: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("add new stock res: %s", res.GetMessage())
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: fmt.Sprintf("i[%s]: new stock %s added", h.instanceName, ticker),
	})
}

func (h *TinyStockExchangeHandler) HandleListStocks(w http.ResponseWriter, r *http.Request) {
	timeoutCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	stream, err := h.tseClient.ListStocks(timeoutCtx, &tseProto.ListStocksRequest{})
	if err != nil {
		log.Errorf("list stocks: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var stocks []*tseProto.Stock
	for {
		stock, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("get stock from stream: %s", err)
			pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		stocks = append(stocks, stock)
	}

	stocksJson, err := json.Marshal(stocks)
	if err != nil {
		log.Errorf("marshal stocks: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("sending %d stocks to client", len(stocks))
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: string(stocksJson),
	})
}

func (h *TinyStockExchangeHandler) HandleNewValueDelta(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Errorf("add new value delta failed, parse form error: %s", err)
		http.Error(w, "parse form error", http.StatusInternalServerError)
		return
	}

	ticker := r.Form.Get("ticker")
	if ticker == "" {
		http.Error(w, "error, ticker empty", http.StatusBadRequest)
		return
	}

	timestampParam := r.Form.Get("ts")
	if timestampParam == "" {
		http.Error(w, "error, timestamp empty", http.StatusBadRequest)
		return
	}
	timestamp, err := strconv.ParseInt(timestampParam, 10, 64)
	if err != nil {
		http.Error(w, "error, timestamp invalid", http.StatusBadRequest)
		return
	}

	deltaParam := r.Form.Get("delta")
	if deltaParam == "" {
		http.Error(w, "error, delta empty", http.StatusBadRequest)
		return
	}
	delta, err := strconv.ParseInt(deltaParam, 10, 64)
	if err != nil {
		http.Error(w, "error, value delta invalid", http.StatusBadRequest)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.tseClient.NewValueDelta(timeoutCtx, &tseProto.StockValueDelta{
		Ticker:    ticker,
		Timestamp: timestamp,
		Delta:     delta,
	})
	if err != nil {
		log.Errorf("add value delta: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("add new stock res: %s", res.GetMessage())
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: fmt.Sprintf("i[%s]: new value delta %d for %s added", h.instanceName, delta, ticker),
	})
}

func (h *TinyStockExchangeHandler) HandleListValueDeltas(w http.ResponseWriter, r *http.Request) {
	timeoutCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	stream, err := h.tseClient.ListStockValueDeltas(timeoutCtx, &tseProto.ListStockValueDeltasRequest{})
	if err != nil {
		log.Errorf("list value deltas: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var deltas []*tseProto.StockValueDelta
	for {
		delta, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("get value delta from stream: %s", err)
			pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		deltas = append(deltas, delta)
	}

	deltasJson, err := json.Marshal(deltas)
	if err != nil {
		log.Errorf("marshal value deltas: %s", err)
		pkg.WriteErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("sending %d value deltas to client", len(deltasJson))
	pkg.WriteJsonResponse(w, http.StatusOK, pkg.ApiResponse{
		Result:  "ok",
		Message: string(deltasJson),
	})
}
