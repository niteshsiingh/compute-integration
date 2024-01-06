package controller

import (
	"context"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/niteshsiingh/matrice-assignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetInstances(c *gin.Context) {
	instance, err := listComputeInstances(project) //lists all available instances
	if err != nil {
		log.Fatal(err)
	}
	var filteredInstances []model.StoreInstance
	for _, a := range instance {
		filter := bson.M{"instanceId": a.Id}
		existingInstance := collection.FindOne(context.TODO(), filter)
		if a.Status != "RUNNING" && existingInstance.Err() == mongo.ErrNoDocuments {
			//checking if the status of instance is not running and this is instance is also not in database then fetch it into the response
			i1, i2, _, i4, i5, i6, i7, i8 := GetInstanceDetails(a)
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
			filteredInstances = append(filteredInstances, b)
		} else if a.Status != "RUNNING" && existingInstance.Err() != mongo.ErrNoDocuments {
			//if it is in database then update the that entry with the new ones
			update := bson.M{
				"$set": bson.M{
					"status": a.Status,
				},
			}
			//update the database
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Fatal(err)
			}
			var ei model.StoreInstance
			_ = collection.FindOne(context.TODO(), filter).Decode(&ei)
			filteredInstances = append(filteredInstances, ei) //fetch and send it to the response
		}

	}
	c.IndentedJSON(http.StatusOK, filteredInstances)
}
