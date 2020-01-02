package prometheus

import (
	"testing"
)

func TestPromQL(t *testing.T) {
    var prom Prometheus
    prom.Address="http://localhost:9090"
    //prom.TCPOutgoing("default","analysis")
    prom.TCPInComing("default","kafka")


}
