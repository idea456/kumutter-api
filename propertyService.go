package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyService struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewPropertyService(ctx context.Context) *PropertyService {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env filed found")
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("test")

	return &PropertyService{
		client: client,
		db:     db,
	}
}

// 1. filter properties by nearest location
// 2. filter listings (from the filtered nearest properties) by their prices
// 3. query shortest route
// 4. query commute cost (if applicable)
func (s *PropertyService) GetPropertiesNear(ctx context.Context, coor []float64) {
	var res bson.M
	query := bson.D{{"location", bson.D{{
		"$near", bson.D{{
			"$geometry", bson.D{{
				"type", "Point",
				"coordinates", coor,
			},
				"$maxDistance", 1000,
			},
		}}}}}}
	coll := s.db.Collection("properties")

	cur, err := coll.Find(ctx, query)
}

func (s *PropertyService) GetPropertiesPrice() {}

func (s *PropertyService) GetPropertiesByDistrict() {}
