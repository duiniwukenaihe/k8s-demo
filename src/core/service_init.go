package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sync"
)

type SVCMapStruct struct {
	data sync.Map
}

func (svcmap *SVCMapStruct) Add(dep *corev1.Service) {
	if list, ok := svcmap.data.Load(dep.Namespace); ok {
		list = append(list.([]*corev1.Service), dep)
		svcmap.data.Store(dep.Namespace, list)
	} else {
		svcmap.data.Store(dep.Namespace, []*corev1.Service{dep})
	}
}
func (svcmap *SVCMapStruct) Update(dep *corev1.Service) error {
	if list, ok := svcmap.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*corev1.Service) {
			if range_dep.Name == dep.Name {
				list.([]*corev1.Service)[i] = dep
				svcmap.data.Store(dep.Namespace, list)
			}
		}
		return nil
	}
	return fmt.Errorf("Service-%s not found", dep.Name)
}

// 删除
func (svcmap *SVCMapStruct) Delete(dep *corev1.Service) {
	if list, ok := svcmap.data.Load(dep.Namespace); ok {
		for i, range_dep := range list.([]*corev1.Service) {
			if range_dep.Name == dep.Name {
				newList := append(list.([]*corev1.Service)[:i], list.([]*corev1.Service)[i+1:]...)
				svcmap.data.Store(dep.Namespace, newList)
				break
			}
		}
	}
}
func (svcmap *SVCMapStruct) ListByNS(ns string) ([]*corev1.Service, error) {
	if list, ok := svcmap.data.Load(ns); ok {
		return list.([]*corev1.Service), nil
	}
	return nil, fmt.Errorf("record not found")
}
func (svcmap *SVCMapStruct) GetService(ns string, name string) (*corev1.Service, error) {
	depList, err := svcmap.ListByNS(ns)
	if err != nil {
		return nil, fmt.Errorf("Service not found")
	}
	for _, item := range depList {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Service not found")
}

var SVCMap *SVCMapStruct

func init() {
	SVCMap = &SVCMapStruct{}
}

type SvcHandler struct {
}

func (s *SvcHandler) OnAdd(obj interface{}) {
	//fmt.Println(obj.(*corev1.Service).Name)
	SVCMap.Add(obj.(*corev1.Service))
}
func (s *SvcHandler) OnUpdate(oldObj, newObj interface{}) {
	err := SVCMap.Update(newObj.(*corev1.Service))
	if err != nil {
		log.Println(err)
	}
}
func (s *SvcHandler) OnDelete(obj interface{}) {
	if s, ok := obj.(*corev1.Service); ok {
		SVCMap.Delete(s)
	}
}
