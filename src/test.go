package main

import (
	"context"
	"fmt"
	"k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func main() {
	dep, _ := lib.K8sClient.AppsV1().Deployments("default").
		Get(context.Background(), "nginx", metav1.GetOptions{})
	selector, err := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(selector.String())
	listOpt := metav1.ListOptions{LabelSelector: selector.String()}
	rs, _ := lib.K8sClient.AppsV1().ReplicaSets("default").
		List(context.Background(), listOpt)
	fmt.Println(dep.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])
	for _, item := range rs.Items {
		fmt.Println(item.Name)
		fmt.Println(item.OwnerReferences)
		fmt.Println(IsCurrentRs(dep, item))
		//s, _ := v1.LabelSelectorAsSelector(item.Spec.Selector)
		fmt.Println(item.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])

	}
}
func IsCurrentRs(dep *v1.Deployment, set v1.ReplicaSet) bool {
	if set.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != dep.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}
	for _, ref := range set.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}

	}
	return false
}

//import (
//	"fmt"
//	"k8s-demo1/src/lib"
//	v1 "k8s.io/api/apps/v1"
//	"k8s.io/apimachinery/pkg/fields"
//	"k8s.io/apimachinery/pkg/util/wait"
//	"k8s.io/client-go/tools/cache"
//)
//
//type DepHandler struct {
//}
//
//func (d *DepHandler) OnAdd(obj interface{}) {}
//func (d *DepHandler) OnUpdate(oldObj, newObj interface{}) {
//	if dep, ok := newObj.(*v1.Deployment); ok {
//		fmt.Println(dep.Name)
//	}
//}
//func (d *DepHandler) OnDelete(obj interface{}) {
//}
//func main() {
//	s, c := cache.NewInformer(cache.NewListWatchFromClient(lib.K8sClient.AppsV1().RESTClient(),
//		"deployments", "default", fields.Everything()),
//		&v1.Deployment{},
//		0,
//		&DepHandler{},
//	)
//	c.Run(wait.NeverStop)
//	s.List()
//}
