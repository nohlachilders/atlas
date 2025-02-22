package server

import (
	"encoding/json"
	"net/http"
)

func (cfg *Config) makeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", healthReponseHandlerFunc)

	mux.Handle("POST /users", CreateUserHandler{
		cfg: cfg,
	})
}

func healthReponseHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		Error string `json:"error"`
	}
	thisError := errorStruct{
		Error: msg,
	}
	res, _ := json.Marshal(thisError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
