package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	svc := NewRoutesService("")
	svc = NewLoggerService(svc)

	mux := NewApiServer(svc)

	routes, err := svc.GetRoutes(context.Background())
	if err != nil {
		log.Fatal(nil)
	}

	fmt.Println(routes)

	log.Fatal(mux.Listen())
}
