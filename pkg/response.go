package pkg

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ApiResponse struct {
	Result  string `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}

func (resp ApiResponse) toJson() []byte {
	apiRespJson, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("failed to marshal response [%v]: %s", resp, err)
		return []byte("{}")
	}
	return apiRespJson
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, apiResp ApiResponse) {
	WriteResponseBytes(w, "application/json", statusCode, apiResp.toJson())
}

func WriteOKJsonResponse(w http.ResponseWriter, statusCode int) {
	apiResp := ApiResponse{Result: "success"}
	WriteResponseBytes(w, "application/json", statusCode, apiResp.toJson())
}

func WriteErrorJsonResponse(w http.ResponseWriter, statusCode int, message string) {
	apiResp := ApiResponse{
		Result:  "failed",
		Message: message,
	}
	WriteResponseBytes(w, "application/json", statusCode, apiResp.toJson())
}

func WritePlainTextResponse(w http.ResponseWriter, statusCode int, message []byte) {
	WriteResponseBytes(w, "text/plain", statusCode, message)
}

func WriteResponseBytes(w http.ResponseWriter, contentType string, statusCode int, message []byte) {
	w.WriteHeader(statusCode)
	if contentType != "" {
		w.Header().Add("Content-Type", contentType)
	}
	if _, err := w.Write(message); err != nil {
		log.Errorf("failed to write response [%s]: %s", message, err)
	}
}
