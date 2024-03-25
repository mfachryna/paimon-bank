package response

import (
	"encoding/json"
	"net/http"

	dto "github.com/mfachryna/paimon-bank/internal/domain/dto/meta"
)

type Response struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ResponseWithMeta struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Meta       dto.Meta    `json:"meta"`
}

var (
	ValidationAndParseBodyError = "Required fields are missing or invalid"
)

func (res *Response) GenerateResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.HttpStatus)
	json.NewEncoder(w).Encode(res)
}
func (res *ResponseWithMeta) GenerateResponseMeta(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.HttpStatus)
	json.NewEncoder(w).Encode(res)
}
