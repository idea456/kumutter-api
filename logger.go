package main

import (
	"context"
	"fmt"
	"time"
)

type LoggerService struct {
	next Service
}

func NewLoggerService(next Service) Service {
	return &LoggerService{
		next: next,
	}
}

func (s *LoggerService) GetRoutes(ctx context.Context) (routes *Route, err error) {
	defer func(start time.Time) {
		fmt.Printf("[%v] routes=%v %v", time.Since(start), routes, err)
	}(time.Now())

	return s.next.GetRoutes(ctx)
}
