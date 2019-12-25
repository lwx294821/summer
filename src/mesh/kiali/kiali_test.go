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
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
     var globalCfg ="D:/golang/workspace/worker/src/factory/workshop/subject/config.yaml"
	 err:=os.Setenv("KIALI_CONF","D:/golang/workspace/worker/src/factory/workshop/subject/mesh/kiali/api/v1/kiali_api_urls.yaml")
     if err !=nil{
     	os.Exit(1)
	 }
	config := viper.New()
	paths, fileName := filepath.Split(globalCfg)
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	if config.IsSet("Kiali") {
		client, err := api.NewClient(api.Config{
			Address:     config.GetString("Kiali.host"),
			UserName:     config.GetString("Kiali.username"),
			Password:     config.GetString("Kiali.password"),
			RoundTripper: nil,
		})
		if err != nil {
			log.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		a:=v1.NewAPI(client)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var p = Path("serviceDetails")
		var args =map[string]string{"namespace":"default","service":"kafka"}
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


}
