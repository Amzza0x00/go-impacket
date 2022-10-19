package main

import (
	"flag"
	"fmt"
	"go-impacket/pkg"
	"go-impacket/pkg/common"
	DCERPCv5 "go-impacket/pkg/dcerpc/v5"
	"go-impacket/pkg/util"
	"log"
	"os"
	"sync"
)

var (
	ip     string
	thread int
	debug  bool
)

func init() {
	flag.StringVar(&ip, "ip", "172.20.10.*", "目标ip或ip段")
	flag.IntVar(&thread, "t", 2000, "线程数量")
	flag.BoolVar(&debug, "debug", false, "开启调试信息")
	flag.Parse()
	fmt.Println(pkg.BANNER)
	if flag.NFlag() < 1 {
		log.Fatalln("Usage: oxidfind -ip 172.20.10.*")
	}
}

func main() {
	ips, err := util.IpParse(ip)
	if err != nil {
		fmt.Printf("[-] ip parse error [%s]: %s\n", ip, err)
		os.Exit(0)
	}
	var wg sync.WaitGroup
	c := make(chan struct{}, thread)
	for _, i := range ips {
		options := common.ClientOptions{
			Host: i,
			Port: 135,
		}
		wg.Add(1)
		go func(ip string) {
			c <- struct{}{}
			defer wg.Done()
			session, err := DCERPCv5.NewTCPSession(options, debug)
			if err != nil {
				fmt.Printf("[-] Connect failed [%s]: %s\n", ip, err)
				return
			}
			rpc, _ := DCERPCv5.TCPTransport()
			rpc.Client = session.Client
			address, err := rpc.ServerAlive2Request(1)
			if err != nil {
				fmt.Println("[-]", err)
				return
			}
			fmt.Printf("[*] %s is alive\n", ip)
			for _, i := range address {
				if i != "" {
					fmt.Printf("[+] NetworkAddr: %s\n", i)
				}
			}
			<-c
		}(i)
	}
	wg.Wait()
}
