package docker

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Container struct {
	Name string
	Id string
	NetWorkMode string
}

/**
   前提安装dockerveth命令工具 docker run --rm -d -v /usr/local/bin:/target dockerveth:latest
 */

//获取指定目录下的镜像tar文件，运行docker load -i [filename].tar
func LoadImages(pwd string){
      //系统路径分隔符pSeparator:=string(os.PathSeparator)
      filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
         if !info.IsDir() && strings.HasSuffix(info.Name(),".tar"){
			var cmd = fmt.Sprintf("docker load -i %s",path)
			execShell(cmd)
		 }
      	return nil
	  })
}

func execShell(command string){
	cmd := exec.Command("bash", "-c", command)
	cmd.Start()
	log.Println(cmd.Args)
	err:=cmd.Wait()
	if err !=nil{
		log.Println(err)
	}
}

/**
  app指组件名称，比如Istio,Kubernetes,需要与配置文件对应上否则会报错。
 */
func PushRegistry(inventory string){
	config := viper.New()
	paths, fileName := filepath.Split(inventory)
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	var app = config.GetString("application")
	if config.IsSet(app) {
		var m = config.GetStringMap(app+".images")
		var hub = config.Get(app+".hub")
		for k,v :=range m{
			if e,ok:=v.(map[string]interface{});ok{
				var s = fmt.Sprintf("%s/%s:%s",e["repo"],k,e["tag"])
				var t =fmt.Sprintf("%s/%s/%s:%s",hub,e["repo"],k,e["tag"])
				var tag = fmt.Sprintf("docker tag %s %s",s,t)
				execShell(tag)
				var push = fmt.Sprintf("docker push %s",t)
				execShell(push)
			}
		}
	}
}

func PullRegistry(inventory string){
	config := viper.New()
	paths, fileName := filepath.Split(inventory)
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	var app = config.GetString("application")
	if config.IsSet(app) {
		var m = config.GetStringMap(app+".images")
		for k,v :=range m{
			if e,ok:=v.(map[string]interface{});ok{
				var s = fmt.Sprintf("%s/%s:%s",e["repo"],k,e["tag"])
				var tag = fmt.Sprintf("docker pull %s ",s)
				execShell(tag)
			}
		}
	}
}


func SaveImages(inventory string){
	config := viper.New()
	paths, fileName := filepath.Split(inventory)
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	var app = config.GetString("application")
	if config.IsSet(app) {
		var m = config.GetStringMap(app+".images")
		for k,v :=range m{
			if e,ok:=v.(map[string]interface{});ok{
				var s = fmt.Sprintf("%s/%s:%s",e["repo"],k,e["tag"])
				var tag = fmt.Sprintf("docker save -o %s.%s.tar %s ",k,e["tag"],s)
				execShell(tag)
			}
		}
	}
}




