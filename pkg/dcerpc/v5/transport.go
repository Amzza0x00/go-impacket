package v5

import (
	"go-impacket/pkg/common"
	"go-impacket/pkg/smb/smb2"
)

type SMBClient struct {
	smb2.Client
}

type TCPClient struct {
	common.Client
}

// 连接封装
// ncacn_np协议的实现
func SMBTransport() (client *SMBClient, err error) {
	return &SMBClient{}, nil
}

// tcp连接封装
func TCPTransport() (client *TCPClient, err error) {
	//session, err := smb2.NewSession(options, debug)
	return &TCPClient{}, nil
}
