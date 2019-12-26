package kubernetes

import "testing"

func TestAPI(t *testing.T) {
	var httpAPI=httpAPI{cfg:Config{"http://localhost:8080","/root/.kube/config"}}
	httpAPI.Nodes()
	httpAPI.Namespaces()
}