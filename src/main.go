package main

import (
	"bufio"
	"context"
	"flag"
	"os"
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
	ctx, cancel := context.WithTimeout(context.Background(),time.Duration(200)*time.Second)
    defer func() {
    	cancel()
	}()
	packet.StartNetSniff("192.168.176.128",ctx)
	in := bufio.NewReader(os.Stdin)
	_, _, err := in.ReadLine()
	if err != nil {
		os.Exit(1)
	}
}
