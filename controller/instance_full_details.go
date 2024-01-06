package controller

import (
	"context"
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"

	"google.golang.org/api/compute/v1"
)

func GetInstanceDetails(instance *compute.Instance) (string, string, int, int, int, int, float64, string) {
	// Extract CPU, GPU, and other details from the instance struct
	cpuPlatform := instance.CpuPlatform
	gpuType := "N/A"
	gpuCount := 0

	if len(instance.GuestAccelerators) > 0 {
		// Sum up the total GPU count across all GPUs attached to the instance
		for _, accelerator := range instance.GuestAccelerators {
			gpuCount += int(accelerator.AcceleratorCount)
			gpuType = accelerator.AcceleratorType
		}
	}
	cpuCount, memoryMB := extractMachineDetails(instance.MachineType)
	//fmt.Println(instance.Disks[0].Source)
	diskSizeGB, err := getDiskSize(project, zone, instance.Disks[0].Source)
	if err != nil {
		log.Fatalf("Error getting disk size: %v", err)
	}
	pricingInfo := 1
	//pricingInfo = instance.Scheduling.MinNodeCpus / 100
	//pricingInfo, err := priceInfo(instance)

	launchTime := instance.CreationTimestamp
	return cpuPlatform, gpuType, gpuCount, cpuCount, memoryMB, diskSizeGB, float64(pricingInfo), launchTime
}

func getDiskSize(project, zone, diskURL string) (int, error) {
	// Create a new Compute Engine API client
	if diskURL == "" {
		diskURL = diskURLcreate
	}
	diskName := path.Base(diskURL)
	service, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("Error creating Compute Engine service: %v", err)
		return 0, err
	}

	// Call the Compute Engine API to get the disk details
	disk, err := service.Disks.Get(project, zone, diskName).Do()
	fmt.Println(disk)
	if err != nil {
		log.Fatalf("Error getting disk details: %v", err)
		return 0, err
	}

	// Extract the disk size
	sizeGB := int(disk.SizeGb)
	return sizeGB, nil
}

func extractMachineDetails(machineTypeURL string) (int, int) {
	// Extract machine type name from the URL
	re := regexp.MustCompile(`projects/[^/]+/zones/[^/]+/machineTypes/(.+)`)
	matches := re.FindStringSubmatch(machineTypeURL)

	if len(matches) == 2 {
		machineTypeName := matches[1]

		// Extract CPU and memory details from the machine type name
		re := regexp.MustCompile(`([0-9]+)cpus-([0-9]+)mb`)
		submatches := re.FindStringSubmatch(machineTypeName)

		if len(submatches) == 3 {
			cpuCount, _ := strconv.Atoi(submatches[1])
			memoryMB, _ := strconv.Atoi(submatches[2])
			return cpuCount, memoryMB
		}
	}

	return 0, 0
}
