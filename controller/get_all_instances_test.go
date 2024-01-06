package controller

import (
	"reflect"
	"testing"

	model "github.com/niteshsiingh/matrice-assignment/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}
func TestGetAllInstances(t *testing.T) {
	a := model.StoreInstance{
		ID:           primitive.NewObjectID(),
		InstanceId:   188,
		InstanceType: "e2-micro",
		Details: model.Instance{
			CPU_Type: "Intel",
			GPU_Type: "N/A",
			Storage:  20,
			Pricing:  2,
		},
		Status: "launched",
	}
	result := insertInstance(a)
	assert(t, result, "e2-micro")
}
