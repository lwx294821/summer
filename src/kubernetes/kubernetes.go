package kubernetes

import (
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Address string
	KubeConfig string
}
type API interface {
	Pods(namespace string)
	Nodes()
	Namespaces()
	Service(namespace,service string)
	Deployment(namespace,app string)
	StatefulSet(namespace string)
	DaemonSet(namespace string)
}

type HttpAPI struct{
     Cfg Config
     API
}
type K8SClientSet struct {
	cli *kubernetes.Clientset
}

func (h *HttpAPI)Client()(*kubernetes.Clientset,error){
	var clientset *kubernetes.Clientset
	var kubeconfig *string
	if h.Cfg.KubeConfig == "" {
		var home = homeDir()
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}else {
		kubeconfig=&h.Cfg.KubeConfig
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil,err
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil,err
	}
	return clientset,nil
}

//查询命名空间中所有Pods
func (h *HttpAPI)Pods(namespace string) {
	clientset,err:= h.Client()
	if clientset == nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		os.Exit(1)
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		fmt.Println(pods)
	}
}


func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}


func (h *HttpAPI)Nodes() *v1.NodeList  {
    clientset,err:= h.Client()
    if clientset == nil{
    	log.Println("Connect kubernetes Fail！！！")
    	log.Println(err)
		return nil
	}
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err !=nil{
		return nil
	}
	return nodes
}

func (h *HttpAPI)Namespaces(){
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		return
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	log.Println(namespaces)
}

func (h *HttpAPI)Service(namespace,service string){
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		return
	}
	svc,err:=clientset.CoreV1().Services(namespace).Get(service,metav1.GetOptions{})
	log.Println(svc)
}

func (h *HttpAPI)Deployment(namespace ,app string) *appsv1.Deployment{
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		return nil
	}
	deploy,err:=clientset.AppsV1().Deployments(namespace).Get(app,metav1.GetOptions{})
	if err !=nil{
		log.Println(err.Error())
		return nil
	}
	return deploy
}

func (h *HttpAPI)WorkLoad(namespace string)map[string]interface{}{
	start :=time.Now()
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		return nil
	}
	var wl=make(map[string]interface{})
	pods,err:=clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err ==nil {
		wl["Pod"]=pods.Items
	}
	ss,err:=clientset.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err == nil {
		wl["StatefulSet"]=ss.Items
	}
	deploy,err:=clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err == nil{
		wl["Deployment"]=deploy.Items
	}
	defer func() {
		end:=time.Since(start)
		log.Printf("query workload time cost = %v",end)
	}()
	return wl
}

func podFilter(pods *v1.PodList){
	var p = pods.Items
	for _,v:=range p{
		var or = v.OwnerReferences
		for _,r:=range or{
			var kind=r.Kind
			log.Println(kind)
		}
	}
}

func DaemonSet(){

}




