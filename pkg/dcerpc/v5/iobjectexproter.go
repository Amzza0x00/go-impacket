package v5

import (
	"bytes"
	"encoding/hex"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
)

// 此文件提供IObjectExporter rpc接口

// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-dcom/8ed0ae33-56a1-44b7-979f-5972f0e9416c
// https://pubs.opengroup.org/onlinepubs/9875999899/CHP21CHP.HTM

// RPC Opnum
const (
	ResolveOxid  = 0
	SimplePing   = 1
	ComplexPing  = 2
	ServerAlive  = 3
	ResolveOxid2 = 4
	ServerAlive2 = 5
)

// ServerAlive2请求结构
type ServerAlive2RequestStruct struct {
	MSRPCHeaderStruct
	AllocHint uint32
	ContextId uint16
	Opnum     uint16
}

type ServerAlive2ResponseStruct struct {
	MSRPCHeaderStruct
	AllocHint       uint32
	ContextId       uint16
	CancelCount     uint8
	Reserved        uint8
	VersionMajor    uint16
	VersionMinor    uint16
	Unknown         uint64
	PpdsaOrBindings AddressStruct
	Reserved2       uint64
}

// 解析的地址结构
type AddressStruct struct {
	NumEntries     uint16
	SecurityOffset uint16
	//Buf            []byte
	//stringBindingStruct
	//securityBindingStruct
}

//type AddressStruct struct {
//	NumEntries     uint16
//	SecurityOffset uint16
//	stringBindingStruct
//	securityBindingStruct
//}

// 绑定地址结构
//type stringBindingStruct struct {
//	TowerId     uint16
//	NetworkAddr []byte // 长度不固定
//}
//
//type securityBindingStruct struct {
//	AuthnSvc  uint16
//	AuthzSvc  uint16
//	PrincName uint16
//}

func NewServerAlive2Request() ServerAlive2RequestStruct {
	header := NewMSRPCHeader()
	header.PacketType = PDURequest
	header.PacketFlags = PDUFault
	header.FragLength = 24
	return ServerAlive2RequestStruct{
		MSRPCHeaderStruct: header,
		AllocHint:         0,
		ContextId:         0,
		Opnum:             ServerAlive2,
	}
}

func NewServerAlive2Response() ServerAlive2ResponseStruct {
	return ServerAlive2ResponseStruct{}
}

//type TcpClient struct {
//	TCPClient
//}

func (c *TCPClient) ServerAlive2Request(callId uint32) (address []string, err error) {
	err = c.MSRPCBind(ms.IID_IObjectExporter, ms.IID_IObjectExporter_VERSION)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	c.Debug("Sending ServerAlive2 request", nil)
	req := NewServerAlive2Request()
	req.CallId = callId
	buf, err := c.TCPSend(req)
	res := NewServerAlive2Response()
	c.Debug("Unmarshalling ServerAlive2 response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 解析address
	addressBuf := buf[40:]
	// 去除securityBinding数据，只保留网卡、ip信息
	securityBindingIndex := bytes.Index(addressBuf, []byte{9, 00})
	newAddressBuf := bytes.Split(addressBuf[:securityBindingIndex], []byte{07, 00})
	for _, i := range newAddressBuf {
		if i != nil {
			address = append(address, string(i))
		}
	}
	return address, nil
}
