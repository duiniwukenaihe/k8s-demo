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
	r.POST("/namespace", service.CreateNameSpace)
	r.POST("/namespace/update", service.UpdateNameSpace)
	r.GET("/deployments", service.ListDeployment)
	r.GET("/service", service.ListService)
	//r.GET("/deployment", service.GetDeployment)
	r.GET("pods", service.ListallPod)
	r.POST("/pod", service.CreatePod)
	r.POST("/pod/delete", service.DeletePod)
	r.POST("/pod/update", service.UpdatePod)
	r.POST("/deployment", service.CreateDep)
	r.POST("/deployment/delete", service.DeleteDep)
	r.POST("/namespace/delete", service.DeleteNameSpace)
	//r.GET("deployment1", service.GetDeployment1)
	core.InitDeployment()
	r.Run()
}
