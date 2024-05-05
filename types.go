package main

import (
	"context"
	"time"
)

type Service interface {
	GetRoutes(ctx context.Context, from Coordinates, to Coordinates) (*Itineary, error)
	GetListingsNear(ctx context.Context, body GetListingsNearBody) (*PropertyServiceResponse[GetListingsNearResponse], error)
}
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CoordinatesBody struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Property struct {
	District    string      `json:"district"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
	rentalRange string      `json:"rentalRange"`
	Facilities  []string    `json:"facilities"`
}

type Amenities struct {
	IsStudio  bool `json:"isStudio"`
	Bathrooms int  `json:"bathrooms"`
	Bedrooms  int  `json:"bedrooms"`
}

type Listing struct {
	Amenities   Amenities `json:"amenities"`
	Price       int       `json:"price"`
	PSF         string    `json:"psf"`
	Area        string    `json:"area"`
	ListingId   string    `json:"listingId"`
	Link        string    `json:"link"`
	IsFurnished string    `json:"furnished"`
	Property    Property  `json:"property"`
	Type        string    `json:"type"`
}

type PlanDataResponse struct {
	Itineraries []Itineary `json:"itineraries"`
}

type PlanResponse struct {
	Plan PlanDataResponse `json:"plan"`
}

type LegStart struct {
	ScheduledTime string `json:"scheduledTime"`
	Estimated     string `json:"estimated"`
}

type LegEnd struct {
	ScheduledTime string `json:"scheduledTime"`
	Estimated     string `json:"estimated"`
}

type LegGeometry struct {
	Length int    `json:"length"`
	Points string `json:"points"`
}

type Leg struct {
	Start LegStart `json:"start"`
	End   LegEnd   `json:"end"`
	Mode  string   `json:"mode"`
	// Duration    int         `json:"duration"`
	LegGeometry LegGeometry `json:"legGeometry"`
	Distance    float64     `json:"distance"`
}

type Itineary struct {
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Duration     int       `json:"duration"`
	WalkDistance float64   `json:"walkDistance"`
	WaitingTime  int       `json:"waitingTime"`
	Legs         []Leg     `json:"legs"`
}
