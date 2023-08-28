package main

import (
	"flag"
	"fmt"
	"github.com/Amzza0x00/go-impacket/pkg"
	"github.com/Amzza0x00/go-impacket/pkg/common"
	DCERPCv5 "github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5"
	"log"
)

var (
	ip    string
	debug bool
)

func init() {
	flag.StringVar(&ip, "ip", "16.16.16.227", "目标ip")
	flag.BoolVar(&debug, "debug", true, "开启调试信息")
	flag.Parse()
	fmt.Println(pkg.BANNER)
	if flag.NFlag() < 1 {
		log.Fatalln("Usage: rpcdump -ip 16.16.16.227")
	}
}

func main() {
	options := common.ClientOptions{
		Host: ip,
		Port: 135,
	}
	session, err := DCERPCv5.NewTCPSession(options, debug)
	if err != nil {
		fmt.Printf("[-] Connect failed [%s]: %s\n", ip, err)
		return
	}
	rpc, _ := DCERPCv5.TCPTransport()
	rpc.Client = session.Client
	err = rpc.RpcBindEpmapper(1)
	if err != nil {
		return
	}
	_, err = rpc.EPMLookupRequest(1)
	if err != nil {
		return
	}
}
