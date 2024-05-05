package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyService struct {
	client *mongo.Client
	db     *mongo.Database
	gql    *graphql.Client
}

type GetListingsNearBody struct {
	MinPrice    int             `json:"min_price"`
	MaxPrice    int             `json:"max_price"`
	MaxDistance int             `json:"max_distance"`
	From        CoordinatesBody `json:"from"`
}

type GetListingsNearResponse struct {
	Name     string    `json:"_id"`
	Property Property  `json:"property"`
	Listings []Listing `json:"items"`
	Itineary Itineary  `json:"itineary"`
}

type PropertyServiceResponse[T any] struct {
	Count int
	Data  []T
}

func NewPropertyService(ctx context.Context) *PropertyService {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env filed found")
	}

	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
	}
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetBSONOptions(bsonOpts)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("test")

	gqlUrl := os.Getenv("ROUTER_URL")
	gql := graphql.NewClient(gqlUrl)

	return &PropertyService{
		client: client,
		db:     db,
		gql:    gql,
	}
}

// 1. filter properties by nearest location
// 2. filter listings (from the filtered nearest properties) by their prices
// 3. query shortest route
// 4. query commute cost (if applicable)
func (s *PropertyService) GetListingsNear(ctx context.Context, body GetListingsNearBody) (*PropertyServiceResponse[GetListingsNearResponse], error) {
	stages := bson.A{
		bson.D{
			{"$geoNear",
				bson.D{
					{"near",
						bson.D{
							{"type", "Point"},
							{"coordinates",
								bson.A{
									body.From.Longitude,
									body.From.Latitude,
								},
							},
						},
					},
					{"distanceField", "string"},
					{"maxDistance", body.MaxDistance},
					{"query", bson.D{}},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"price",
						bson.D{
							{"$gt", body.MinPrice},
							{"$lt", body.MaxPrice},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$property.name"},
					{"items",
						bson.D{
							{"$addToSet",
								bson.D{
									{"link", "$link"},
									{"listingId", "$listingId"},
									{"price", "$price"},
									{"amenities", "$amenities"},
									{"psf", "$psf"},
									{"area", "$area"},
									{"property", "$property"},
								},
							},
						},
					},
				},
			},
		},
	}

	coll := s.db.Collection("listings")
	cursor, err := coll.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}

	var results []GetListingsNearResponse
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i, result := range results {
		results[i].Property = result.Listings[0].Property
	}

	return &PropertyServiceResponse[GetListingsNearResponse]{
		Count: len(results),
		Data:  results,
	}, nil
}

func (s *PropertyService) GetRoutes(ctx context.Context, from Coordinates, to Coordinates) (*Itineary, error) {
	query := fmt.Sprintf(`
	query {
		plan(
		  from: {lat: %f, lon: %f},
		  to: {lat: %f, lon: %f}
		) {
		  itineraries {
			start
			end
			duration
			walkDistance
			waitingTime
			legs {
			  start {
				scheduledTime
				estimated {
				  time
				}
			  }
			  end {
				scheduledTime
				estimated {
				  time
				}
			  }
			  mode
			  duration
			  legGeometry {
				length
				points
			  }
			  distance
			}
		  }
		}
	  }
	`, from.Latitude, from.Longitude, to.Latitude, to.Longitude)

	fmt.Println(from, to)
	req := graphql.NewRequest(query)

	var plan PlanResponse
	err := s.gql.Run(ctx, req, &plan)
	if err != nil {
		log.Fatal(err)
	}

	if len(plan.Plan.Itineraries) == 0 {
		return &Itineary{}, nil
	}
	return &plan.Plan.Itineraries[0], nil
}
