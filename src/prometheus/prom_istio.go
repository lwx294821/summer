package prometheus

import (
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
	"os"
	"time"
)

var address = "http://localhost:9090"

type Prometheus struct {
	Address string
}


const IstioTcpSentBytesTotal = "istio_tcp_sent_bytes_total"
// 统计源流出流量
func (p *Prometheus)TCPOutgoing(namespace,source string ) {
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
		"reporter":"source",
		"source_workload_namespace":namespace,
		"source_workload":source,
	}
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
    var query =p.PromQL(IstioTcpSentBytesTotal,match)
	log.Println("irate("+query+"[5m])")
	result, warnings, err := api.QueryRange(ctx, "irate("+query+"[5m])", r)
	if err != nil {
		log.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		log.Printf("Warnings: %v\n", warnings)
	}
	log.Println(result)
}


func (p *Prometheus)PromQL(metric string,labels map[string]string) string{
     var buff bytes.Buffer
     buff.WriteString(metric)
     buff.WriteString("{")
     for k,v:=range labels{
     	buff.WriteString(k+"=\""+v+"\",")
	 }
	 buff.WriteString("}")
	 return buff.String()
}

const IstioTcpReceivedBytesTotal = "istio_tcp_received_bytes_total"
// 统计目的工作负载进入流量
func (p *Prometheus)TCPInComing(namespace,destination string){
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
		"reporter":"destination",
		"destination_workload_namespace":namespace,
		"destination_workload":destination,
	}
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	var query =p.PromQL(IstioTcpReceivedBytesTotal,match)
	log.Println("irate("+query+"[5m])")
	result, warnings, err := api.QueryRange(ctx,"irate("+query+"[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Println("Result:")
	fmt.Printf("Result:\n%v\n", result)
}

func Workload(){

}







