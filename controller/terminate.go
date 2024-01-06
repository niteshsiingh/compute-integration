package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/niteshsiingh/matrice-assignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func terminate(instance_id uint64) {
	hexID := strconv.FormatUint(instance_id, 10)
	// fmt.Printf("Type of idd: %T\n", hexID)
	// fmt.Println(hexID)

	var instance model.StoreInstance
	filter := bson.M{"instanceId": instance_id}
	err := collection.FindOne(context.Background(), filter).Decode(&instance)
	if err != nil {
		log.Fatal(err)
		return
	}
	launchTime, err := time.Parse(time.RFC3339, instance.LaunchTime)
	if err != nil {
		log.Fatal(err)
		return
	}
	currentTime := time.Now()
	usedHours := currentTime.Sub(launchTime).Minutes() / 60
	effectiveCost := usedHours * instance.Details.Pricing
	update := bson.M{"$set": bson.M{"status": "terminated", "cost": effectiveCost}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	defer deleteInstance(hexID)
	if err != nil {
		log.Fatal(err)
	}

}

func TerminateInstance(c *gin.Context) {
	id := c.Query("instance_id")
	idd, _ := strconv.ParseUint(id, 10, 64)
	terminate(idd)
	var filteredInstances []model.StoreInstance
	filter := bson.M{"instanceId": id}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var instance model.StoreInstance
		err := cursor.Decode(&instance)
		if err != nil {
			log.Fatal(err)
		}
		if instance.InstanceId == idd {
			filteredInstances = append(filteredInstances, instance)
			c.IndentedJSON(http.StatusOK, filteredInstances)
			return
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(filteredInstances)
}

func deleteInstance(instanceId string) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("/Users/niteshsingh/VS Code/Go/alert-brook-410312-a19f83b314c5.json"))
	if err != nil {
		log.Fatalf("Failed to create Compute Engine API client: %v", err)
	}

	_ = fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, instanceId)

	op, err := computeService.Instances.Delete(project, zone, instanceId).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to delete instance: %v", err)
	}

	err = waitForOperation(computeService, op.Name)
	if err != nil {
		log.Fatalf("Failed waiting for operation: %v", err)
	}

	fmt.Printf("Instance %s terminated successfully\n", instanceId)
}
