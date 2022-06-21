package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	"k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type Deployment struct {
	Namespace           string `json:"namespace"`
	Name                string `json:"name"`
	Replicas            int32  `json:"replicas"`
	AvailableReplicas   int32  `json:"available-replicas"`
	UnavailableReplicas int32  `json:"unavailable-replicas"`
	Images              string `json:"images"`
	Ports               string `json:"ports"`
	CreateTime          string `json:"CreateTime"`
	Labels              string `json:"labels"`
	Pods                []*Pod `json:"pods"`
}

func ListDeployment(g *gin.Context) {
	ns := g.Query("ns")
	deplist, _ := core.DepMap.ListByNS(ns)
	//dps, err := K8sClient.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
	//if err != nil {
	//	g.Error(err)
	//}
	ret := make([]*Deployment, 0)
	for _, item := range deplist {
		ret = append(ret, &Deployment{
			Namespace:           item.Namespace,
			Name:                item.Name,
			Replicas:            item.Status.Replicas,
			AvailableReplicas:   item.Status.AvailableReplicas,
			UnavailableReplicas: item.Status.UnavailableReplicas,
			Images:              item.Spec.Template.Spec.Containers[0].Image,
			//Labels:              item.GetLabels(),
			Pods:       GetPodsByDep(*item),
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:03:04"),
		})

	}
	g.JSON(200, ret)
	return
}
func Createdep(dep Deployment) (*v1.Deployment, error) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: dep.Namespace,
			Labels:    dep.GetLabels(),
		},
		Spec: v1.DeploymentSpec{
			Replicas: &dep.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: dep.GetSelectors(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   dep.Name,
					Labels: dep.GetSelectors(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  dep.GetImageName(),
							Image: dep.Images,
							Ports: dep.GetPorts(),
							//Ports: []corev1.ContainerPort{
							//{
							//	Name:          "web",
							//	Protocol:      corev1.ProtocolTCP,
							//	ContainerPort: 80,
							//},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	newdep, err := lib.K8sClient.AppsV1().Deployments(dep.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return newdep, nil
}

func CreateDep(g *gin.Context) {
	var newDep Deployment
	if err := g.ShouldBind(&newDep); err != nil {
		g.JSON(500, err)
	}
	newdep, err := Createdep(newDep)
	if err != nil {
		g.JSON(500, err)
	}
	newDep1 := Deployment{
		Namespace:  newdep.Namespace,
		Name:       newdep.Name,
		Pods:       GetPodsByDep(*newdep),
		CreateTime: newdep.CreationTimestamp.Format("2006-01-02 15:03:04"),
	}
	g.JSON(200, newDep1)
}

func GetLabels(m map[string]string) string {
	labels := ""
	// aa=xxx,xxx=xx

	for k, v := range m {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", k, v)
	}
	return labels
}
func GetPodsByDep(dep v1.Deployment) []*Pod {
	rsLabelsMap, err := core.RSMap.GetRsLabelsByDeployment(&dep)
	//fmt.Println(rsLabelsMap)
	if err != nil {
		log.Fatal(err)
	}
	pods, err := core.PodMap.ListByRsLabels(dep.Namespace, rsLabelsMap)
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]*Pod, 0)
	for _, pod := range pods {
		//
		if core.RSMap.GetRsLabelsByDeploymentname(&dep) == pod.OwnerReferences[0].Name {
			ret = append(ret, &Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Images:    pod.Spec.Containers[0].Image,
				NodeName:  pod.Spec.NodeName,
				Labels:    pod.Labels,
				Status:    string(pod.Status.Phase),
				//IsReady:   GetPodIsReady(*pod),
				//	Message:    GetPodMessage(*pod),
				//Message:      core.EventMap.GetMessage(pod.Namespace, "Pod", pod.Name),
				//HostIp:       pod.Status.HostIP,
				//PodIp:        pod.Status.PodIP,
				//RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
				CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			})
		}
	}
	return ret
}

//func GetDeployment(ns string,name string) {
//	ret := make([]*Deployment, 0)
//	ret = append(ret, &Deployment{
//		Namespace:           ret.Namespace,
//		Name:                dps.Name,
//		Replicas:            dps.Status.Replicas,
//		AvailableReplicas:   dps.Status.AvailableReplicas,
//		UnavailableReplicas: dps.Status.UnavailableReplicas,
//		Images:              dps.Spec.Template.Spec.Containers[0].Image,
//		CreateTime:          dps.CreationTimestamp.Format("2006-01-02 15:03:04"),
//		Labels:              dps.Labels,
//		Pods:                GetPodsByDep(ns, dps),
//	})
//	g.JSON(200, ret)
//	return
//}

//func GetPodsByDep(ns string, dep *v1.Deployment) []*Pod {
//	ctx := context.Background()
//	listopt := metav1.ListOptions{
//		//LabelSelector: GetLabels(dep.Spec.Selector.MatchLabels),
//		LabelSelector: common.GetRsLableByDeployment(dep),
//	}
//	list, err := K8sClient.CoreV1().Pods(ns).List(ctx, listopt)
//	if err != nil {
//		panic(err.Error())
//	}
//	pods := make([]*Pod, len(list.Items))
//	for i, pod := range list.Items {
//		pods[i] = &Pod{
//			Namespace:  pod.Namespace,
//			Name:       pod.Name, //获取 pod名称
//			Status:     string(pod.Status.Phase),
//			Images:     pod.Spec.Containers[0].Image,
//			NodeName:   pod.Spec.NodeName, //所属节点
//			Labels:     pod.Labels,
//			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"), //创建时间
//		}
//	}
//
//	return pods
//
//}
