package utils

import (
	"encoding/json"
	"net/http"

	"github.com/berkayaydmr/language-learning-api/pkg/customerr"
)

func RespondWithError(w http.ResponseWriter, err error) {
	resp, code := customerr.NewErrorResponse(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func Respond(w http.ResponseWriter, data interface{}, code int) {
	resp, err := json.Marshal(data)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
