package kiali

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"summer/src/mesh/kiali/api"
	v1 "summer/src/mesh/kiali/api/v1"
	"time"
)





func Query(method string,args map[string]string,cfg api.Config){
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}
	a:=v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var p = Path(method)
	result, warnings, err :=a.Get(ctx,p,args,time.Now())
	if err != nil {
		log.Printf("Error querying Kiali: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		log.Printf("Warnings: %v\n", warnings)
	}
	log.Printf("Result:\n%v\n", result)
}


func Path(method string) string{
	config := viper.New()
	KialiConf := os.Getenv("KIALI_CONF")
	paths, fileName := filepath.Split(KialiConf)
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	if config.IsSet(method) {
		return config.GetString(method)
	}
	return "/"
}
