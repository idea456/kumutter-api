package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) Listen() error {
	http.HandleFunc("/listings", s.handleGetPropertiesNear)
	return http.ListenAndServe(":6969", nil)
}

func (s *ApiServer) handleGetPropertiesNear(w http.ResponseWriter, r *http.Request) {
	var body GetListingsNearBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	listingsRes, err := s.svc.GetListingsNear(context.Background(), body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d properties!\n", listingsRes.Count)

	var wg sync.WaitGroup
	for i, listing := range listingsRes.Data {
		wg.Add(1)
		go func(res GetListingsNearResponse, index int) {
			itineary, err := s.svc.GetRoutes(context.Background(), res.Listings[0].Property.Coordinates, Coordinates(body.From))
			if err != nil {
				panic(err)
			}

			listingsRes.Data[index].Itineary = *itineary
			wg.Done()
		}(listing, i)
	}

	wg.Wait()

	writeJSON(w, http.StatusOK, listingsRes)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
