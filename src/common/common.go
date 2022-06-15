package common

import (
	"context"
	"fmt"
	"k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetRsLableByDeployment(dep *v1.Deployment) string {
	selector, _ := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	listOpt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	rs, _ := lib.K8sClient.AppsV1().ReplicaSets(dep.Namespace).List(context.Background(), listOpt)
	for _, item := range rs.Items {
		if IsCurrentRsByDep(dep, &item) {
			s, err := metav1.LabelSelectorAsSelector(item.Spec.Selector)
			if err != nil {
				return ""
			}
			return s.String()
		}
	}
	return ""
}
func GetRsLabelByDeployment_ListWatch(dep *v1.Deployment, rslist []*v1.ReplicaSet) (map[string]string, error) {
	for _, item := range rslist {
		if IsCurrentRsByDep(dep, item) {
			return item.Labels, nil
		}
	}
	return nil, fmt.Errorf("find replicaset error")
}

func IsRsFromDep(dep *v1.Deployment, set v1.ReplicaSet) bool {
	for _, ref := range set.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}
	}
	return false

}
func IsCurrentRsByDep(dep *v1.Deployment, set *v1.ReplicaSet) bool {
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
