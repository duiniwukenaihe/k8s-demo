package service

import (
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	corev1 "k8s.io/api/core/v1"
)

type Service struct {
	Name       string
	Type       string
	ClusterIp  string
	ExternalIp []string
	Ports      []corev1.ServicePort
	Select     map[string]string
}

func ListService(g *gin.Context) {
	ns := g.Query("ns")
	//svc, err := K8sClient.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{})
	svclist, _ := core.SVCMap.ListByNS(ns)
	ret := make([]*Service, 0)
	for _, item := range svclist {
		ret = append(ret, &Service{
			Name:       item.Name,
			Type:       string(item.Spec.Type),
			ClusterIp:  item.Spec.ClusterIP,
			ExternalIp: item.Spec.ExternalIPs,
			Ports:      item.Spec.Ports,
			Select:     item.Spec.Selector,
		})

	}
	g.JSON(200, ret)
	return
}
