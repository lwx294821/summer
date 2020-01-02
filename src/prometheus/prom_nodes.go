package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
	"os"
	"time"
)





//单个主机上每个网络接口的每秒上传流量
const networkReceiced ="node_network_receive_bytes_total"
func (p *Prometheus)NetReceived(instance string){
	client, err := api.NewClient(api.Config{
		Address: p.Address,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}
	api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var match = map[string]string{
		"job":"node-exporter",
		"instance":instance,
	}
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	var query =p.PromQL(networkReceiced,match)
	log.Println("rate("+query+"[5m])")
	result, warnings, err := api.QueryRange(ctx, "rate("+query+"[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Node Nerwork: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Println("Result:")
	fmt.Printf("Result:\n%v\n", result)
}
