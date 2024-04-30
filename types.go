package main

import "context"

type Route struct {
	Legs     []interface{}
	Shape    string
	Duration string
	Distance int
}

type Service interface {
	GetRoutes(context.Context) (*Route, error)
}
