package main

import (
	"flag"
	"summer/src/docker"
)

var kind string
var username string
var password string
var timeout string

func init() {
	flag.StringVar(&kind, "kind", "kubernetes", "Specify the client type to be used.")
	flag.StringVar(&username, "username", "admin", "Login Username.")
	flag.StringVar(&password, "password", "admin", "Login Password.")
	flag.StringVar(&timeout, "timeout", "5m", "Login Password.")
}

func main() {
   docker.PushRegistry("","Istio")
	//var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary
	//b, err := json_iterator.Marshal(dev)

}

