package main

import (
	"flag"
	"fmt"
	"go-impacket/pkg"
	"go-impacket/pkg/common"
	DCERPCv5 "go-impacket/pkg/dcerpc/v5"
	"go-impacket/pkg/smb/smb2"
	"go-impacket/pkg/util"
	"log"
	"os"
)

// 1.查找可用共享目录
// 2.上传文件
// 3.打开远程服务
// 4.创建服务并启动

var (
	user     string
	domain   string
	password string
	hash     string
	target   string
	port     int
	file     string
	path     string
	debug    bool
	service  string
)

func init() {
	flag.StringVar(&user, "user", "", "用户名,默认为空")
	flag.StringVar(&domain, "domain", "de1ay", "用户名,默认为de1ay")
	flag.StringVar(&password, "pass", "", "密码,默认为空")
	flag.StringVar(&hash, "hash", "", "哈希,默认为空")
	flag.StringVar(&target, "target", "", "目标地址,默认为空")
	flag.IntVar(&port, "port", 445, "目标端口,默认为445")
	flag.StringVar(&file, "file", "", "要安装的服务可执行文件,默认为空")
	flag.StringVar(&path, "path", "", "可执行文件的目录路径,默认为空")
	flag.BoolVar(&debug, "debug", false, "开启调试信息,默认为关闭")
	flag.StringVar(&service, "service", "", "创建的服务名称,默认为随机4位字符")
	flag.Parse()
	fmt.Println(pkg.BANNER)
	if flag.NFlag() < 5 {
		log.Fatalln("Usage: psexec -target 172.20.10.2 -user administrator -hash 32ed87bdb5fdc5e9cba88547376818d4 -file test.exe -path ./test/")
	}
	if target == "" {
		log.Fatalln("目标地址为空")
	}
}

func main() {
	options := common.ClientOptions{
		Host:     target,
		Port:     port,
		Domain:   domain,
		User:     user,
		Password: password,
		Hash:     hash,
	}
	session, err := smb2.NewSession(options, debug)
	if err != nil {
		fmt.Printf("[-] Login failed [%s]: %s\n", target, err)
		os.Exit(0)
	}
	defer session.Close()
	if session.IsAuthenticated {
		fmt.Printf("[+] Login successful [%s]\n", target)
	}
	var serviceName string
	if service == "" {
		serviceName = string(util.Random(4))
	} else {
		serviceName = service
	}
	rpc, _ := DCERPCv5.SMBTransport()
	rpc.Client = *session
	// 创建服务并启动
	servicename, _, _ := rpc.ServiceInstall(serviceName, file, path)
	fmt.Printf("[+] Service name is [%s]\n", servicename)
}
