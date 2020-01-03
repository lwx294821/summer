package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"summer/src/local"
	"summer/src/local/packet"
	"time"
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
	dev := local.FindAllNetWorkDevs()
	var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json_iterator.Marshal(dev)
	if err == nil {
		fmt.Println(string(b))

	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200)*time.Second)
	defer func() {
		cancel()
	}()
	packet.StartNetSniff("192.168.176.128", ctx)
	in := bufio.NewReader(os.Stdin)
	_, _, err = in.ReadLine()
	if err != nil {
		os.Exit(1)

	}
}
