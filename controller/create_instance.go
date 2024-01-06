package controller

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/niteshsiingh/matrice-assignment/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func CreateInstance(c *gin.Context) {
	var instanceTyp model.InstanceData
	if err := c.BindJSON(&instanceTyp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := createInstance(instanceTyp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "creation done"})
}
func createInstance(instanceTyp model.InstanceData) error {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)
	randomString := strconv.Itoa(randomNumber)
	instanceType := instanceTyp.Type
	instanceName := "nitesh-" + instanceType + "-" + randomString
	// instanceName := instanceTyp.Name
	// fmt.Println(instanceName)

	computeService, err := createComputeEngine()
	if err != nil {
		log.Fatal(err)
	}
	statusValue := "launched"
	instanc := &compute.Instance{
		Name:        instanceName,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", zone, instanceType),
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Name: "default",
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
			},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "STATUS",
					Value: &statusValue,
				},
			},
		},
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,         // The first disk must be a boot disk.
				AutoDelete: true,         //Optional
				Mode:       "READ_WRITE", //Mode should be READ_WRITE or READ_ONLY
				Interface:  "SCSI",       //SCSI or NVME - NVME only for SSDs
				InitializeParams: &compute.AttachedDiskInitializeParams{

					DiskName:    "worker-instance-boot-disk" + randomString,
					SourceImage: "projects/centos-cloud/global/images/family/centos-7",
					DiskType:    fmt.Sprintf("projects/%s/zones/%s/diskTypes/pd-ssd", project, zone),
					DiskSizeGb:  20,
				},
			},
		},
	}

	op, err := computeService.Instances.Insert(project, zone, instanc).Do()
	if err != nil {
		fmt.Println(err)
	}
	var firewallRule *compute.Firewall
	firewallRuleName := "allow-9000-9100"

	firewallRules, err := computeService.Firewalls.List(project).Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, existingRule := range firewallRules.Items {
		if existingRule.Name == firewallRuleName {
			firewallRule = existingRule
			break
		}
	}

	if firewallRule == nil {
		firewallRule = &compute.Firewall{
			Name:        generateUniqueName(firewallRuleName),
			Description: "Allow traffic on ports 9000-9100",
			Allowed: []*compute.FirewallAllowed{
				{
					IPProtocol: "tcp",
					Ports:      []string{"900-910"},
				},
			},
			SourceRanges: []string{"0.0.0.0/0"},
		}
		_, err := computeService.Firewalls.Insert(project, firewallRule).Do()
		if err != nil {
			if isAlreadyExistsError(err) {
				log.Fatal("retry with a different name and error is: ", err)
			} else {
				log.Fatal(err)
			}
		}

	}
	err = waitForOperation(computeService, op.Name)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	createdInstance, err := computeService.Instances.Get(project, zone, instanc.Name).Do()
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	if len(createdInstance.Disks) > 0 {
		diskURLcreate = createdInstance.Disks[0].Source
	}

	_, i2, _, i4, i5, i6, _, _ := GetInstanceDetails(instanc)
	b := model.StoreInstance{
		ID:           primitive.NewObjectIDFromTimestamp(time.Now()),
		InstanceId:   createdInstance.Id,
		InstanceType: instanceType,
		Details: model.Instance{
			CPU_Type:  createdInstance.CpuPlatform,
			GPU_Type:  i2,
			GPU_Count: i4,
			Memory:    i5,
			Storage:   i6,
			Pricing:   1,
		},
		LaunchTime: createdInstance.CreationTimestamp,
		Status:     "launched",
	}
	insertInstance(b)
	return nil
}

func generateUniqueName(baseName string) string {
	timestamp := time.Now().UnixNano()
	randomSuffix := rand.Intn(10000)
	uniqueName := fmt.Sprintf("%s-%d-%04d", baseName, timestamp, randomSuffix)
	return uniqueName
}

func isAlreadyExistsError(err error) bool {
	apiErr, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}
	return apiErr.Code == 409 && strings.Contains(apiErr.Message, "alreadyExists")
}
