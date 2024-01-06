package controller

import (
	"context"
	"fmt"
	"log"

	model "github.com/niteshsiingh/matrice-assignment/models"
)

func insertInstance(instance model.StoreInstance) string {
	_, err := collection.InsertOne(context.Background(), instance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted instance with type: ", instance.InstanceType)
	return instance.InstanceType
}
