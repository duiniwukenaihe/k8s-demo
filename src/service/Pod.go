package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	. "k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Images      string            `json:"images"`
	NodeName    string            `json:"nodename"`
	CreateTime  string            `json:"createtime"`
	Annotations map[string]string `json:"annotations"`
	Ports       string            `json:"ports"`
	//IsReady    bool
	//Message      string
	//HostIp       string
	//PodIp        string
	//RestartCount int32
	Labels map[string]string `json:"labels"`
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
			Namespace: item.Namespace,
			Name:      item.Name,
			Status:    string(item.Status.Phase),
			Labels:    item.Labels,
			NodeName:  item.Spec.NodeName,
			Images:    item.Spec.Containers[0].Image,
			//IsReady:   GetPodIsReady(*item),
			//Message: GetPodMessage(*item),
			//Message:      core.EventMap.GetMessage(item.Namespace, "Pod", item.Name),
			//HostIp:       item.Status.HostIP,
			//PodIp:        item.Status.PodIP,
			//RestartCount: item.Status.ContainerStatuses[0].RestartCount,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})

	}
	g.JSON(200, ret)
	return
}
func Createpod(pod Pod) (*corev1.Pod, error) {
	newpod, err := K8sClient.CoreV1().Pods(pod.Namespace).Create(context.TODO(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: pod.Name, Image: pod.Images},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return newpod, err
}
func CreatePod(g *gin.Context) {
	var NewPod Pod
	if err := g.ShouldBind(&NewPod); err != nil {
		g.JSON(500, err)
	}
	pod, err := Createpod(NewPod)
	if err != nil {
		g.JSON(500, err)
	}
	newpod := Pod{
		Namespace:   pod.Namespace,
		Name:        pod.Name,
		Images:      pod.Spec.Containers[0].Image,
		CreateTime:  pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Annotations: nil,
	}
	g.JSON(200, newpod)
}
func updatepod(pod Pod) (*corev1.Pod, error) {
	newpod, err := K8sClient.CoreV1().Pods(pod.Namespace).Update(context.TODO(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: pod.Name, Image: pod.Images, Ports: pod.GetPorts()},
			},
		},
	}, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return newpod, err
}
func UpdatePod(g *gin.Context) {
	var NewPod Pod
	if err := g.ShouldBind(&NewPod); err != nil {
		g.JSON(500, err)
	}
	pod, err := updatepod(NewPod)
	if err != nil {
		g.JSON(500, err)
	}
	newpod := Pod{
		Namespace:   pod.Namespace,
		Name:        pod.Name,
		Images:      pod.Spec.Containers[0].Image,
		CreateTime:  pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Annotations: pod.ObjectMeta.Annotations,
	}
	g.JSON(200, newpod)
}
func DeletePod(g *gin.Context) {
	var NewPod Pod
	if err := g.ShouldBind(&NewPod); err != nil {
		g.JSON(500, err)
	}
	err := K8sClient.CoreV1().Pods(NewPod.Namespace).Delete(context.TODO(), NewPod.Name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	}
	g.JSON(200, "ok")
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
			//Message:      core.EventMap.GetMessage(item.Namespace, "Pod", item.Name),
			//RestartCount: item.Status.ContainerStatuses[0].RestartCount,
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
