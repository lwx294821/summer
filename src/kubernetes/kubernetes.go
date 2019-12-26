package kubernetes

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
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
	App(namespace,app string)
}

type httpAPI struct{
     cfg Config
     API
}
type K8SClientSet struct {
	cli *kubernetes.Clientset
}

func (h *httpAPI)Client()(*kubernetes.Clientset,error){
	var clientset *kubernetes.Clientset
	var kubeconfig *string
	if h.cfg.KubeConfig == "" {
		var home = homeDir()
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}else {
		kubeconfig=&h.cfg.KubeConfig
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

func (h *httpAPI)Pods(namespace string) {
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


func (h *httpAPI)Nodes(){
    clientset,err:= h.Client()
    if clientset == nil{
    	log.Println("Connect kubernetes Fail！！！")
    	log.Println(err)
		return
	}
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	log.Println(nodes)
}

func (h *httpAPI)Namespaces(){
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		return
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	log.Println(namespaces)
}

func (h *httpAPI)Service(namespace,service string){
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		return
	}
	svc,err:=clientset.CoreV1().Services(namespace).Get(service,metav1.GetOptions{})
	log.Println(svc)
}

func (h *httpAPI)App(namespace ,app string){
	clientset,err:= h.Client()
	if clientset ==nil{
		log.Println("Connect kubernetes Fail！！！")
		log.Println(err)
		return
	}
	deploy,err:=clientset.AppsV1().Deployments(namespace).Get(app,metav1.GetOptions{})
	log.Println(deploy)

}