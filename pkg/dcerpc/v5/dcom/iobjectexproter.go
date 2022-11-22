package dcom

import (
	//"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5"
	"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5/rpcrt"
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
	rpcrt.MSRPCHeaderStruct
	AllocHint uint32
	ContextId uint16
	Opnum     uint16
}

type ServerAlive2ResponseStruct struct {
	rpcrt.MSRPCHeaderStruct
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
	header := rpcrt.NewMSRPCHeader()
	header.PacketType = rpcrt.PDURequest
	header.PacketFlags = rpcrt.PDUFault
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
