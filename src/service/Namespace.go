package service

import (
	"context"
	"github.com/gin-gonic/gin"
	. "k8s-demo1/src/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type Time struct {
	time.Time `protobuf:"-"`
}
type Namespace struct {
	Name       string
	CreateTime Time `json:"CreateTime"`
	Status     string
	Labels     map[string]string
}

func ListNamespace(g *gin.Context) {
	ns, err := K8sClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		g.Error(err)
		return
	}
	ret := make([]*Namespace, 0)
	for _, item := range ns.Items {
		ret = append(ret, &Namespace{
			Name:       item.Name,
			CreateTime: Time(item.CreationTimestamp),
			Status:     string(item.Status.Phase),
			Labels:     item.Labels,
		})

	}
	g.JSON(200, ret)
	return
}
