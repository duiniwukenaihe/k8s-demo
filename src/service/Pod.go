package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	. "k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	Namespace  string
	Name       string
	Status     string
	Images     string
	NodeName   string
	CreateTime string
	Message    string
	Labels     map[string]string
}

func ListallPod(g *gin.Context) {
	ns := g.Query("ns")

	//pods, err := K8sClient.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	pods, err := core.PodMap.ListByNS(ns)
	if err != nil {
		g.Error(err)
	}
	ret := make([]*Pod, 0)
	for _, item := range pods {

		ret = append(ret, &Pod{
			Namespace:  item.Namespace,
			Name:       item.Name,
			Status:     string(item.Status.Phase),
			Labels:     item.Labels,
			NodeName:   item.Spec.NodeName,
			Images:     item.Spec.Containers[0].Image,
			Message:    GetPodMessage(*item),
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})

	}
	g.JSON(200, ret)
	return
}
func ListPodsByLabel(ns string, labels []map[string]string) (ret []*Pod) {
	list, err := core.PodMap.ListByRsLabels(ns, labels)
	if err != nil {
		return nil
	}

	for _, item := range list {
		ret = append(ret, &Pod{
			Name:      item.Name,
			Namespace: item.Namespace,
			Images:    item.Spec.Containers[0].Image,
			NodeName:  item.Spec.NodeName,
			Status:    string(item.Status.Phase),
			//Message: GetPodMessage(*item),
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:22:33"),
		})
	}

	return
}

func ListPod(ns string, dep *v1.Deployment) []*Pod {
	//labelSelector := make(map[string]string)

	listopt := metav1.ListOptions{
		LabelSelector: GetLabels(dep.Spec.Selector.MatchLabels),
	}

	podList, err := K8sClient.CoreV1().Pods(ns).List(context.Background(), listopt)
	if err != nil {
		fmt.Println(err.Error())
	}
	pods := make([]*Pod, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = &Pod{
			Namespace:  pod.Namespace,
			Name:       pod.Name, //获取 pod名称
			Status:     string(pod.Status.Phase),
			Images:     pod.Spec.Containers[0].Image,
			NodeName:   pod.Spec.NodeName, //所属节点
			Labels:     pod.Labels,
			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"), //创建时间
		}
	}
	return pods
}
