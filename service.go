package main

import (
	"context"
	"log"

	graphql "github.com/machinebox/graphql"
)

type Service interface {
	GetRoutes(context.Context) (*Route, error)
}

type RoutesService struct {
	client *graphql.Client
}

func NewRoutesService(url string) Service {
	client := graphql.NewClient("https://machinebox.io/graphql")
	return &RoutesService{
		client: client,
	}
}

func (s *RoutesService) GetRoutes(ctx context.Context) (*Route, error) {
	req := graphql.NewRequest(`
	
	`)
	req.Header.Set("Content-Type", "application/json")

	var route Route
	err := s.client.Run(ctx, req, &route)
	if err != nil {
		log.Fatal(err)
	}

	return &route, nil
}
