package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/niteshsiingh/matrice-assignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func getInstanceDetail(instanceType string) error {
	localInstances, err := fetchInfoByMachineType(instanceType)
	if err != nil {
		return err
	}
	for _, a := range localInstances {
		filter := bson.D{{Key: "instanceId", Value: a.Id}}
		i1, i2, _, i4, i5, i6, i7, i8 := GetInstanceDetails(a)
		existingInstance := collection.FindOne(context.TODO(), filter)
		if existingInstance.Err() == mongo.ErrNoDocuments {
			machineType := a.MachineType
			b := model.StoreInstance{
				ID:           primitive.NewObjectIDFromTimestamp(time.Now()),
				InstanceId:   a.Id,
				InstanceType: path.Base(machineType),
				Details: model.Instance{
					CPU_Type:  i1,
					GPU_Type:  i2,
					GPU_Count: i4,
					Memory:    i5,
					Storage:   i6,
					Pricing:   i7,
				},
				LaunchTime: i8,
				Status:     a.Status,
			}
			insertInstance(b)
		} else {
			update := bson.M{
				"$set": bson.M{
					"details": bson.M{
						"cpu-type":  i1,
						"gpu-type":  i2,
						"gpu-count": i4,
						"memory":    i5,
						"storage":   i6,
						"pricing":   i7,
					},
				},
			}

			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	return nil
}
func fetchInfoByMachineType(machineType string) ([]*compute.Instance, error) {
	ctx := context.Background()

	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		return nil, fmt.Errorf("Error creating Compute Engine service: %v", err)
	}
	machineTypeURL, err := computeService.MachineTypes.Get(project, zone, machineType).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Error getting machine type URL: %v", err)
	}
	instances, err := computeService.Instances.List(project, zone).Filter(fmt.Sprintf("machineType eq '%s'", machineTypeURL.SelfLink)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Error listing instances: %v", err)
	}
	return instances.Items, nil
}

func GetInstanceDetail(c *gin.Context) {
	id := c.Query("instance_type")
	result := getInstanceDetail(id)
	if result != nil {
		log.Fatal(result)
	}
	filter := bson.M{"type": id}
	var filteredInstances []model.StoreInstance
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
		filteredInstances = append(filteredInstances, instance)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	if len(filteredInstances) == 0 {
		createInstance(model.InstanceData{Type: id})
	}
	c.IndentedJSON(http.StatusOK, filteredInstances)
}
