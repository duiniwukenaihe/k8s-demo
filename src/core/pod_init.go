package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sync"
)

type PodMapStruct struct {
	data sync.Map
}

func (podmap *PodMapStruct) Add(pod *corev1.Pod) {
	key := pod.Namespace
	if value, ok := podmap.data.Load(key); ok {
		value = append(value.([]*corev1.Pod), pod)
		podmap.data.Store(key, value)
	} else {
		podmap.data.Store(key, []*corev1.Pod{pod})
	}
}
func (podmap *PodMapStruct) Update(pod *corev1.Pod) error {
	key := pod.Namespace
	if value, ok := podmap.data.Load(key); ok {

		for index, p := range value.([]*corev1.Pod) {
			if p.Name == pod.Name {
				value.([]*corev1.Pod)[index] = pod
				podmap.data.Store(key, value)
				return nil
			}
		}
	}

	return fmt.Errorf("pod-%s not found", pod.Name)
}
func (podmap *PodMapStruct) Delete(pod *corev1.Pod) {
	key := pod.Namespace
	if value, ok := podmap.data.Load(key); ok {
		for index, p := range value.([]*corev1.Pod) {
			if p.Name == pod.Name {
				value = append(value.([]*corev1.Pod)[0:index], value.([]*corev1.Pod)[index+1:]...)
				podmap.data.Store(key, value)
				return
			}
		}
	}
}
func (podmap *PodMapStruct) ListByNS(ns string) ([]*corev1.Pod, error) {

	if ns != "" {
		if list, ok := podmap.data.Load(ns); ok {
			return list.([]*corev1.Pod), nil
		}
	}
	return nil, fmt.Errorf("pods not found")
}

func (podmap *PodMapStruct) ListByRsLabels(ns string, labels []map[string]string) ([]*corev1.Pod, error) {
	pods, err := podmap.ListByNS(ns)
	if err != nil {
		return nil, err
	}
	ret := make([]*corev1.Pod, 0)
	for _, pod := range pods {
		for _, label := range labels {
			if IsValidLabel(pod.Labels, label) {
				ret = append(ret, pod)
			}
		}
	}
	return ret, nil
}

type PodHandler struct {
}

var PodMap *PodMapStruct

func init() {
	PodMap = &PodMapStruct{}
}
func (podmap *PodHandler) OnAdd(obj interface{}) {
	PodMap.Add(obj.(*corev1.Pod))
}
func (podmap *PodHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
}
func (podmap *PodHandler) OnDelete(obj interface{}) {
	PodMap.Delete(obj.(*corev1.Pod))
}
