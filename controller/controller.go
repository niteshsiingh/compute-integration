package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

const connectionString = "YOUR_CONNECTION_STRING"
const dbName = "Your_DATABASE_NAME"
const colName = "Your_COLLECTION_NAME"
const zone = "Your_DESIRED_ZONE"
const project = "PPROJECT-ID"
const portRangeStart = 9000
const portRangeEnd = 9100
const API_KEY = "YOUR API-KEY"

var collection *mongo.Collection
var diskURLcreate string

func init() {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection established")
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance ready")
}

func listComputeInstances(project string) ([]*compute.Instance, error) {
	ctx := context.Background()

	service, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Compute Engine service: %v", err)
	}

	instances, err := service.Instances.List(project, zone).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Error listing instances: %v", err)
		return nil, err
	}
	return instances.Items, nil
}

func createComputeEngine() (*compute.Service, error) {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Error creating Compute Engine service: %v", err)
		return nil, err
	}
	return service, err
}
