package v5

import (
	"fmt"
	"github.com/Amzza0x00/go-impacket/pkg/common"
	"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5/rpcrt"
	"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5/tcprpct"
	"net"
)

//type SMBClient struct {
//	smb2.Client
//}

//type TCPClient struct {
//	common.Client
//}
type SMBClient = rpcrt.SMBClient
type TCPClient = tcprpct.TCPClient

// 连接封装
// ncacn_np协议的实现
func SMBTransport() (client *SMBClient, err error) {
	return &SMBClient{}, nil
}

// tcp连接封装
func NewTCPSession(opt common.ClientOptions, debug bool) (client *TCPClient, err error) {
	address := fmt.Sprintf("%s:%d", opt.Host, opt.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	client = &TCPClient{}
	client.WithOptions(&opt)
	client.WithConn(conn)
	client.WithDebug(debug)
	return client, nil
}

func TCPTransport() (client *TCPClient, err error) {
	return &TCPClient{}, nil
}
