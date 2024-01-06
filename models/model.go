package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instance struct {
	CPU_Type  string  `json:"cpu-type,omitempty" bson:"cpu-type,omitempty"`
	GPU_Type  string  `json:"gpu-type,omitempty" bson:"gpu-type,omitempty"`
	GPU_Count int     `json:"gpu-count,omitempty" bson:"gpu-count,omitempty"`
	Memory    int     `json:"memory,omitempty" bson:"memory,omitempty"`
	Storage   int     `json:"storage,omitempty" bson:"storage,omitempty"`
	Pricing   float64 `json:"pricing,omitempty" bson:"pricing,omitempty"`
}

type StoreInstance struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	InstanceId   uint64             `json:"instanceId,omitempty" bson:"instanceId,omitempty"`
	InstanceType string             `json:"type,omitempty" bson:"type,omitempty"`
	Details      Instance           `json:"details,omitempty" bson:"details,omitempty"`
	LaunchTime   string             `json:"launch-time,omitempty" bson:"launch-time,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	Cost         float64            `json:"cost,omitempty" bson:"cost,omitempty"`
}

type InstanceData struct {
	Type string `json:"type"`
}
