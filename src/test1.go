package main

//func main() {
//}

//import (
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"k8s-demo1/src/lib"
//	v1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/util/wait"
//	"k8s.io/client-go/informers"
//	"log"
//	"sync"
//	"time"
//)
//
//type DeploymentMap struct {
//	data sync.Map
//}
//
//func (depmap *DeploymentMap) Add(dep *v1.Deployment) {
//	if list, ok := depmap.data.Load(dep.Namespace); ok {
//		list = append(list.([]*v1.Deployment), dep)
//		depmap.data.Store(dep.Namespace, list)
//	} else {
//		depmap.data.Store(dep.Namespace, []*v1.Deployment{dep})
//	}
//}
//
//var DepMap *DeploymentMap
//
//func init() {
//	DepMap = &DeploymentMap{}
//}
//
//type PodMap struct {
//	data sync.Map
//}
//type PodHandler struct {
//}
//
//func (p *PodHandler) OnAdd(obj interface{}) {
//	fmt.Println(obj.(*corev1.Pod).Name)
//
//}
//func (p *PodHandler) OnUpdate(oldObj, newObj interface{}) {
//	if pods, ok := newObj.(*corev1.Pod); ok {
//		fmt.Println(pods.Name)
//	}
//}
//func (p *PodHandler) OnDelete(obj interface{}) {
//}
//
//type DepHandler struct {
//}
//
//func (d *DepHandler) OnAdd(obj interface{}) {
//	//fmt.Println(obj.(*v1.Deployment).Name)
//	DepMap.Add(obj.(*v1.Deployment))
//}
//func (d *DepHandler) OnUpdate(oldObj, newObj interface{}) {
//	if dep, ok := newObj.(*v1.Deployment); ok {
//		fmt.Println(dep.Name)
//	}
//}
//func (d *DepHandler) OnDelete(obj interface{}) {
//}
//func main() {
//	factory := informers.NewSharedInformerFactory(lib.K8sClient, 0)
//	podinformer := factory.Core().V1().Pods()
//	podinformer.Informer().AddEventHandler(&PodHandler{})
//	depinformer := factory.Apps().V1().Deployments()
//	depinformer.Informer().AddEventHandler(&DepHandler{})
//	factory.Start(wait.NeverStop)
//	c, _ := context.WithTimeout(context.Background(), time.Second*3)
//	select {
//	case <-c.Done():
//		log.Fatal("time out")
//	default:
//		r := gin.New()
//		r.GET("/", func(c *gin.Context) {
//			var res []string
//			DepMap.data.Range(func(key, value interface{}) bool {
//				if key == "kube-system" {
//					for _, item := range value.([]*v1.Deployment) {
//						res = append(res, item.Name)
//					}
//				}
//				return true
//
//			})
//			c.JSON(200, res)
//
//		})
//		r.Run(":8080")
//	}
//}

//func main() {
//	s, c := cache.NewInformer(cache.NewListWatchFromClient(lib.K8sClient.CoreV1().RESTClient(),
//		"pods", "default", fields.Everything()),
//		&corev1.Pod{},
//		0,
//		&PodHandler{},
//	)
//	c.Run(wait.NeverStop)
//	s.List()
//}

//type DepHandler struct{}
//
//func (d DepHandler) OnAdd(obj interface{}) {}
//
//func (d DepHandler) OnUpdate(oldObj, newObj interface{}) {
//	//TODO implement me
//	if dep, ok := newObj.(*v1.Deployment); ok {
//		fmt.Println(dep.Name)
//	}
//}
//
//func (d DepHandler) OnDelete(obj interface{}) {}

//func (this *DepHandler) OnDelete(obj interface{}) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (this *DepHandler) OnAdd(obj interface{}) {}
//func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
//}
//func main() {
//	s, c := cache.NewInformer(
//		cache.NewListWatchFromClient(lib.K8sClient.AppsV1().RESTClient(),
//			"deployments", "default", fields.Everything()),
//		&v1.Deployment{},
//		0,
//		&DepHandler{},
//	)
//	c.Run(wait.NeverStop)
//	s.List()
//}
//
//import (
//	"fmt"
//	"k8s-demo1/src/lib"
//	"k8s.io/api/apps/v1"
//	"k8s.io/apimachinery/pkg/util/wait"
//	"k8s.io/client-go/informers"
//)
//
//type DepHandler struct{}
//
//func (this *DepHandler) OnAdd(obj interface{}) {}
//func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
//	if dep, ok := newObj.(*v1.Deployment); ok {
//		fmt.Println(dep.Name)
//	}
//}
//
//func (this *DepHandler) OnDelete(obj interface{}) {}
//
//func main() {
//	fact := informers.NewSharedInformerFactory(lib.K8sClient, 0)
//	depInformer := fact.Apps().V1().Deployments()
//	depInformer.Informer().AddEventHandler(&DepHandler{})
//
//	fact.Start(wait.NeverStop)
//	select {}
//
//}
