package core

import (
	"errors"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"log"
	"sync"
)

type RSMapStruct struct {
	Data sync.Map
}

func (rsmap *RSMapStruct) Add(rs *v1.ReplicaSet) {
	key := rs.Namespace
	if value, ok := rsmap.Data.Load(key); ok {
		value = append(value.([]*v1.ReplicaSet), rs)
		rsmap.Data.Store(key, value)
	} else {
		rsmap.Data.Store(key, []*v1.ReplicaSet{rs})
	}
}
func (rsmap *RSMapStruct) Update(rs *v1.ReplicaSet) error {
	key := rs.Namespace
	if value, ok := rsmap.Data.Load(key); ok {
		for index, r := range value.([]*v1.ReplicaSet) {
			if r.Name == rs.Name {
				value.([]*v1.ReplicaSet)[index] = rs
				rsmap.Data.Store(key, value)
				return nil
			}
		}
	}

	return fmt.Errorf("rs-%s not found", rs.Name)
}

func (rsmap *RSMapStruct) Delete(rs *v1.ReplicaSet) {
	key := rs.Namespace
	if value, ok := rsmap.Data.Load(key); ok {
		for index, r := range value.([]*v1.ReplicaSet) {
			if r.Name == rs.Name {
				value = append(value.([]*v1.ReplicaSet)[0:index], value.([]*v1.ReplicaSet)[index+1:]...)
				rsmap.Data.Store(key, value)
				return
			}
		}
	}
}

func (rsmap *RSMapStruct) ListByNS(ns string) ([]*v1.ReplicaSet, error) {
	if list, ok := rsmap.Data.Load(ns); ok {
		return list.([]*v1.ReplicaSet), nil
	}
	return nil, errors.New("rs record not found")
}

func (rsmap *RSMapStruct) GetRsLabelsByDeployment(deploy *v1.Deployment) ([]map[string]string, error) {
	rs, err := rsmap.ListByNS(deploy.Namespace)
	if err != nil {
		return nil, err
	}
	ret := make([]map[string]string, 0)
	for _, item := range rs {
		//if item.Annotations["deployment.kubernetes.io/revision"] != deploy.Annotations["deployment.kubernetes.io/revision"] {
		//	continue
		//}
		for _, v := range item.OwnerReferences {
			if v.Name == deploy.Name {
				ret = append(ret, item.Labels)
				break
			}
		}
	}
	return ret, nil
}
func (rsmap *RSMapStruct) GetRsLabelsByDeploymentname(deploy *v1.Deployment) string {
	rs, err := rsmap.ListByNS(deploy.Namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range rs {
		//if item.Annotations["deployment.kubernetes.io/revision"] != deploy.Annotations["deployment.kubernetes.io/revision"] {
		//	continue
		//}
		for _, v := range item.OwnerReferences {
			if v.Name == deploy.Name {
				return item.Name
			}
		}
	}
	return ""
}

type RSHandler struct {
}

func (rsmap *RSHandler) OnAdd(obj interface{}) {
	RSMap.Add(obj.(*v1.ReplicaSet))
}
func (rsmap *RSHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := RSMap.Update(newObj.(*v1.ReplicaSet))
	if err != nil {
		log.Println(err)
	}
}
func (rsmap *RSHandler) OnDelete(obj interface{}) {
	RSMap.Delete(obj.(*v1.ReplicaSet))
}

var RSMap *RSMapStruct

func init() {
	RSMap = &RSMapStruct{}
}
