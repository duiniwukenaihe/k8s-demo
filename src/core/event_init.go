package core

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"sync"
)

var EventMap *EventMapStruct

type EventMapStruct struct {
	data sync.Map
}

func (eventmap EventMapStruct) GetMessage(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s-%s-%s", ns, kind, name)
	if v, ok := eventmap.data.Load(key); ok {
		return v.(*v1.Event).Message
	}
	return ""
}

type EventHandler struct{}

func (eventmap *EventHandler) storeData(obj interface{}, isDelete bool) {
	if event, ok := obj.(*v1.Event); ok {
		key := fmt.Sprintf("%s-%s-%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isDelete {
			EventMap.data.Store(key, event)
		} else {
			EventMap.data.Delete(key)
		}
	}
}

func (eventmap *EventHandler) OnAdd(obj interface{}) {
	eventmap.storeData(obj, false)
}

func (eventmap *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	eventmap.storeData(newObj, false)
}

func (eventmap *EventHandler) OnDelete(obj interface{}) {
	eventmap.storeData(obj, true)
}

func init() {
	EventMap = &EventMapStruct{}
}
