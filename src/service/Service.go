package service

import (
	"context"
	"github.com/gin-gonic/gin"
	. "k8s-demo1/src/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	Name       string
	Type       string
	ClusterIp  string
	ExternalIp []string
	Ports      []string
	Select     map[string]string
}

func ListService(g *gin.Context) {
	ns := g.Query("ns")
	svc, err := K8sClient.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		g.Error(err)
		return
	}
	ret := make([]*Service, 0)
	for _, item := range svc.Items {
		ret = append(ret, &Service{
			Name:       item.Name,
			Type:       string(item.Spec.Type),
			ClusterIp:  item.Spec.ClusterIP,
			ExternalIp: item.Spec.ExternalIPs,
			Select:     item.Spec.Selector,
		})

	}
	g.JSON(200, ret)
	return
}
