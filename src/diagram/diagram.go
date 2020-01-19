package diagram

import (
	"log"
	"summer/src/diagram/node"
	"summer/src/kubernetes"
)

/**
  实现方法: 调用K8S API查询工作负载的拓扑节点列表，
  workload: Deployment\Pod\StatefulSet
  传参:命名空间和工作负载类型
 */
var httpApi kubernetes.HttpAPI
func init(){
	httpApi=kubernetes.HttpAPI{Cfg:kubernetes.Config{Address: "http://localhost:8080", KubeConfig: "/root/.kube/config"}}
}
//查询工作负载之部署
func GetDeployment(namespace,name string)  {

}




func QueryNodes(namespace,workload,kind string) []node.Node{
	var tnode  []node.Node
	switch kind {
	case "Deployment":
		  tnode,err:=deployment(namespace,workload)
		  if err !=nil{
		  	log.Println(err)
		  }
		  return tnode
	case "Pod":
		return nil
	case "StatefulSet":
		return nil
	default:
		log.Println("Don't find match type !!!")
		return tnode
	}

}
func deployment(namespace,workload string) ([]node.Node,error){
	var d = httpApi.Deployment(namespace,workload)
	var ns  []node.Node
	ns=append(ns,node.Node{
		ID:       string(d.UID),
		Category: d.Kind,
		Text:     d.Name,
		Status:   d.Status.String(),
	})
	return ns,nil
}

//拓扑连接数据
func QueryLinks(namespace string){

}




