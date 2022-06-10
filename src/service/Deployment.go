package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-demo1/src/core"
	v1 "k8s.io/api/apps/v1"
	"log"
)

type Deployment struct {
	Namespace           string
	Name                string
	Replicas            int32
	AvailableReplicas   int32
	UnavailableReplicas int32
	Images              string
	CreateTime          string
	Labels              map[string]string
	Pods                []*Pod
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
			Labels:              item.GetLabels(),
			Pods:                GetPodsByDep(*item),
			CreateTime:          item.CreationTimestamp.Format("2006-01-02 15:03:04"),
		})

	}
	g.JSON(200, ret)
	return
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
	if err != nil {
		log.Fatal(err)
	}

	pods, err := core.PodMap.ListByRsLabels(dep.Namespace, rsLabelsMap)
	if err != nil {
		log.Fatal(err)
	}
	ret := make([]*Pod, 0)

	for _, pod := range pods {
		if core.RSMap.GetRsLabelsByDeploymentname(&dep) == pod.OwnerReferences[0].Name {
			ret = append(ret, &Pod{
				Name:       pod.Name,
				Namespace:  pod.Namespace,
				Images:     pod.Spec.Containers[0].Image,
				NodeName:   pod.Spec.NodeName,
				Labels:     pod.Labels,
				Status:     string(pod.Status.Phase),
				CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			})
		}
	}
	return ret
}

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
