package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{}
}

func (s *ApiServer) Listen() error {
	http.HandleFunc("/routes", s.handleGetRoutes)
	return http.ListenAndServe(":4000", nil)
}

func (s *ApiServer) handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := s.svc.GetRoutes(context.Background())
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
	}

	writeJSON(w, http.StatusOK, routes)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
