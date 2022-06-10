package core

import (
	"fmt"
	"k8s-demo1/src/lib"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"log"
	"sync"
)

type DeploymentMap struct {
	data sync.Map
}

func (depmap *DeploymentMap) Add(dep *v1.Deployment) {
	if list, ok := depmap.data.Load(dep.Namespace); ok {
		list = append(list.([]*v1.Deployment), dep)
		depmap.data.Store(dep.Namespace, list)
	} else {
		depmap.data.Store(dep.Namespace, []*v1.Deployment{dep})
	}
}
func (depmap *DeploymentMap) Update(dep *v1.Deployment) error {
	if list, ok := depmap.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*v1.Deployment) {
			if range_dep.Name == dep.Name {
				list.([]*v1.Deployment)[i] = dep
				depmap.data.Store(dep.Namespace, list)
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found", dep.Name)
}

// 删除
func (depmap *DeploymentMap) Delete(dep *v1.Deployment) {
	if list, ok := depmap.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*v1.Deployment) {
			if range_dep.Name == dep.Name {
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				depmap.data.Store(dep.Namespace, newList)
				break
			}
		}
	}
}
func (depmap *DeploymentMap) ListByNS(ns string) ([]*v1.Deployment, error) {
	if list, ok := depmap.data.Load(ns); ok {
		return list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}

var DepMap *DeploymentMap

func init() {
	DepMap = &DeploymentMap{}
}

type DepHandler struct {
}

func (d *DepHandler) OnAdd(obj interface{}) {
	//fmt.Println(obj.(*v1.Deployment).Name)
	DepMap.Add(obj.(*v1.Deployment))
}
func (d *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	err := DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}
func (d *DepHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*v1.Deployment); ok {
		DepMap.Delete(d)
	}
}

func InitDeployment() {
	factory := informers.NewSharedInformerFactory(lib.K8sClient, 0)
	depinformer := factory.Apps().V1().Deployments()
	depinformer.Informer().AddEventHandler(&DepHandler{})
	Podinformer := factory.Core().V1().Pods()
	Podinformer.Informer().AddEventHandler(&PodHandler{})
	Rsinformer := factory.Apps().V1().ReplicaSets()
	Rsinformer.Informer().AddEventHandler(&RSHandler{})
	factory.Start(wait.NeverStop)
}
