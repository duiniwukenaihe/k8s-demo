package main

import (
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	"k8s-demo1/src/service"
	//	"k8s.io/client-go/informers/core"
)

func main() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, "hello")
	})
	r.GET("/namespaces", service.ListNamespace)
	r.GET("/deployments", service.ListDeployment)
	r.GET("/service", service.ListService)
	//r.GET("/deployment", service.GetDeployment)
	r.GET("pods", service.ListallPod)
	//	r.GET("deployment1", service.GetDeployment1)
	core.InitDeployment()
	r.Run()
}
