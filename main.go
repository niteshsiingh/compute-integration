package main

import (
	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/matrice-assignment/controller"
)

func main() {
	router := gin.Default()
	router.GET("/instances", controller.GetInstances)
	router.PUT("/instance", controller.GetInstanceDetail)
	router.PUT("/instances", controller.TerminateInstance)
	router.POST("/instance", controller.CreateInstance)
	router.Run("localhost:8080")
}
