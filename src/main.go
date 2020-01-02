package main

import (
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"summer/src/local"
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
	dev:=local.FindAllNetWorkDevs()
	var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json_iterator.Marshal(dev)
	if err ==nil{
		fmt.Println(string(b))
	}

}
