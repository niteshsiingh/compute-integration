package controller

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

func waitForOperation(computeService *compute.Service, operationName string) error {
	ctx := context.Background()
	for {
		op, err := computeService.ZoneOperations.Get(project, zone, operationName).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("Failed to get operation status: %v", err)
		}

		if op.Status == "DONE" {
			if op.Error != nil {
				return fmt.Errorf("Operation completed with error: %v", op.Error)
			}
			return nil
		}
		fmt.Println("Waiting for operation to complete...")
		time.Sleep(5 * time.Second)
	}
}
