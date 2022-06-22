package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	. "k8s-demo1/src/lib"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type Time struct {
	time.Time `protobuf:"-"`
}
type Namespace struct {
	Name        string            `json:"name"`
	CreateTime  time.Time         `json:"CreateTime"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
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
			CreateTime: item.CreationTimestamp.Time,
			Status:     string(item.Status.Phase),
			Labels:     item.Labels,
		})

	}
	g.JSON(200, ret)
	return
}
func create(ns Namespace) (*v1.Namespace, error) {
	ctx := context.Background()
	newNamespace, err := K8sClient.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ns.Name,
			Labels: ns.Labels,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return newNamespace, err
}
func updatenamespace(ns Namespace) (*v1.Namespace, error) {
	ctx := context.Background()
	newNamespace, err := K8sClient.CoreV1().Namespaces().Update(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ns.Name,
			Labels: ns.Labels,
		},
	}, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return newNamespace, err
}
func CreateNameSpace(g *gin.Context) {
	var nameSpace Namespace
	if err := g.ShouldBind(&nameSpace); err != nil {
		g.JSON(500, err)
	}
	namespace, err := create(nameSpace)
	if err != nil {
		g.JSON(500, err)
	}
	ns := Namespace{
		Name:        namespace.Name,
		CreateTime:  namespace.CreationTimestamp.Time,
		Status:      string(namespace.Status.Phase),
		Labels:      nil,
		Annotations: nil,
	}
	g.JSON(200, ns)
}
func UpdateNameSpace(g *gin.Context) {
	var nameSpace Namespace
	if err := g.ShouldBind(&nameSpace); err != nil {
		g.JSON(500, err)
	}
	namespace, err := updatenamespace(nameSpace)
	if err != nil {
		g.JSON(500, err)
	}
	ns := Namespace{
		Name:        namespace.Name,
		CreateTime:  namespace.CreationTimestamp.Time,
		Status:      string(namespace.Status.Phase),
		Labels:      namespace.Labels,
		Annotations: nil,
	}
	g.JSON(200, ns)
}
func DeleteNameSpace(g *gin.Context) {
	var nameSpace Namespace
	if err := g.ShouldBind(&nameSpace); err != nil {
		g.JSON(500, err)
	}
	err := K8sClient.CoreV1().Namespaces().Delete(context.Background(), nameSpace.Name, metav1.DeleteOptions{})
	if err != nil {
		g.JSON(500, err)
	}
	g.JSON(200, "Namespace has delete")
}
