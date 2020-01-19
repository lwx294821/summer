package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var kind string
var timeout string

var server *http.Server

func init() {
	flag.StringVar(&kind, "kind", "kubernetes", "Specify the client type to be used.")
	flag.StringVar(&timeout, "timeout", "5m", "Operation timeout.")
}

func main() {


}
func StartServer() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.Handle("/api", &myHandler{})
	mux.HandleFunc("/bye", sayBye)

	server = &http.Server{
		Addr:         ":2020",
		WriteTimeout: time.Second * 4,
		Handler:      mux,
	}

	go func() {
		// 接收退出信号
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()
	err := server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
	log.Fatal("Server exited")

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//get param then make a decision
	
	_,err:=w.Write([]byte("Hello World!"))
	if err!=nil{
		log.Println(err.Error())
	}
}

// 关闭http
func sayBye(w http.ResponseWriter, r *http.Request) {
	_,err:=w.Write([]byte("bye bye!!! Shutdown the server"))
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Println([]byte("shutdown the server err"))
	}
}



