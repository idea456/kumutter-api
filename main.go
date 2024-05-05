package main

import (
	"context"
	"log"
	"sync"
)

func routeWorker(wg *sync.WaitGroup, queue chan int) {

}

func main() {
	// svc := NewRoutesService("http://ec2-13-229-96-163.ap-southeast-1.compute.amazonaws.com:8080/otp/routers/default/index/graphql")
	// svc = NewLoggerService(svc)

	// var trxExchange = CoordinatesBody{
	// 	Latitude:  3.1428187982993845,
	// 	Longitude: 101.7181168996529,
	// }
	propertySvc := NewPropertyService(context.Background())
	mux := NewApiServer(propertySvc)

	// nearestListings := propertySvc.GetListingsNear(context.Background(), trxExchange)

	// propertySvc.GetListingsNear(context.Background(), GetListingsNearBody{
	// 	MinPrice:    1000,
	// 	MaxPrice:    2000,
	// 	MaxDistance: 2000,
	// 	From:        trxExchange,
	// })

	// var wg sync.WaitGroup
	// var itinearies []Itineary
	// for _, listing := range nearestListings {
	// 	wg.Add(1)
	// 	go func(property Listing) {
	// 		fmt.Printf("Requesting route from %s to TRX Exchange...\n", listing.Property.Name)
	// 		defer wg.Done()
	// 		itineary, err := svc.GetRoutes(context.Background(), listing.Property.Coordinates, trxExchange)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Printf("Route from %s to TRX Exchange: %f metres\n", listing.Property.Name, itineary.WalkDistance)
	// 		itinearies = append(itinearies, *itineary)
	// 	}(listing)
	// }

	// wg.Wait()

	// fmt.Println(routes)

	log.Fatal(mux.Listen())
}
