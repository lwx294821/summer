package docker

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
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
func findAllContainerVeth(){

}


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
		log.Println(err.Error())
	}
}

/**
  app指组件名称，比如
 */
func PushRegistry(hub string,app string){
	config := viper.New()
	paths, fileName := filepath.Split("/root/summer/src//config.yaml")
	config.AddConfigPath(paths)
	config.SetConfigName(fileName)
	var suffix = path.Ext(fileName)
	config.SetConfigType(suffix[1:])
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	if config.IsSet(app) {
		var m = config.GetStringMap(app+".images")
		var hub = config.Get(app+".hub")
		

	}

}



