package docker


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

