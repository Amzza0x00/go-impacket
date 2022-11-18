package v5

import (
	"encoding/hex"
	"errors"
	"github.com/Amzza0x00/go-impacket/pkg/encoder"
	"github.com/Amzza0x00/go-impacket/pkg/ms"
	"github.com/Amzza0x00/go-impacket/pkg/smb/smb2"
	"github.com/Amzza0x00/go-impacket/pkg/util"
)

// 此文件提供ms-rpce封装
// DCE/RPC RPC over SMB 协议实现
// https://pubs.opengroup.org/onlinepubs/9629399/

// MSRPC 标准头
type MSRPCHeaderStruct struct {
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32 //4字节，小端排序，0x10
	FragLength         uint16 //2字节，整个结构的长度
	AuthLength         uint16
	CallId             uint32
}

func NewMSRPCHeader() MSRPCHeaderStruct {
	return MSRPCHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         0,
		PacketFlags:        0,
		DataRepresentation: 16,
		AuthLength:         0,
	}
}

type MSRPCRequestHeaderStruct struct {
	MSRPCHeaderStruct
	AllocHint uint32 `smb:"len:Buffer"` //Buffer的长度
	ContextId uint16
	OpNum     uint16
	Buffer    interface{}
}

// 函数绑定请求结构
type MSRPCBindStruct struct {
	MSRPCHeaderStruct
	MaxXmitFrag uint16 //4字节，发送大小协商
	MaxRecvFrag uint16 //4字节，接收大小协商
	AssocGroup  uint32
	NumCtxItems uint8
	Reserved    uint8
	Reserved2   uint16
	CtxItem     CtxEItemStruct
}

// 函数绑定响应结构
type MSRPCBindAckStruct struct {
	MSRPCHeaderStruct
	MaxXmitFrag   uint16
	MaxRecvFrag   uint16
	AssocGroup    uint32
	ScndryAddrlen uint16
	ScndryAddr    []byte `smb:"count:ScndryAddrlen"` //取决管道的长度
	NumResults    uint8
	CtxItem       CtxEItemResponseStruct
}

// PDU CtxItem结构
type CtxEItemStruct struct {
	ContextId      uint16
	NumTransItems  uint8
	Reserved       uint8
	AbstractSyntax SyntaxIDStruct
	TransferSyntax SyntaxIDStruct
}

type SyntaxIDStruct struct {
	UUID    []byte `smb:"fixed:16"`
	Version uint32
}

// PDU CtxItem响应结构
type CtxEItemResponseStruct struct {
	AckResult      uint16
	AckReason      uint16
	TransferSyntax []byte `smb:"fixed:16"` //16字节
	SyntaxVer      uint32
}

// PDU PacketType
// https://pubs.opengroup.org/onlinepubs/9629399/chap12.htm
const (
	PDURequest            = 0
	PDUPing               = 1
	PDUResponse           = 2
	PDUFault              = 3
	PDUWorking            = 4
	PDUNoCall             = 5
	PDUReject             = 6
	PDUAck                = 7
	PDUCl_Cancel          = 8
	PDUFack               = 9
	PDUCancel_Ack         = 10
	PDUBind               = 11
	PDUBind_Ack           = 12
	PDUBind_Nak           = 13
	PDUAlter_Context      = 14
	PDUAlter_Context_Resp = 15
	PDUShutdown           = 17
	PDUCo_Cancel          = 18
	PDUOrphaned           = 19
)

// PDU PacketFlags
// https://pubs.opengroup.org/onlinepubs/9629399/chap12.htm
const (
	PDUFlagReserved_01 = 0x01
	PDUFlagLastFrag    = 0x02
	PDUFlagPending     = 0x03
	PDUFlagFrag        = 0x04
	PDUFlagNoFack      = 0x08
	PDUFlagMayBe       = 0x10
	PDUFlagIdemPotent  = 0x20
	PDUFlagBroadcast   = 0x40
	PDUFlagReserved_80 = 0x80
)

// NDR 传输标准
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rpce/b6090c2b-f44a-47a1-a13b-b82ade0137b2
const (
	NDRSyntax   = "8a885d04-1ceb-11c9-9fe8-08002b104860" //Version 02, NDR64 data representation protocol
	NDR64Syntax = "71710533-BEBA-4937-8319-B5DBEF9CCC36" //Version 01, NDR64 data representation protocol
)

// 函数绑定响应
func NewMSRPCBindAck() MSRPCBindAckStruct {
	return MSRPCBindAckStruct{
		CtxItem: CtxEItemResponseStruct{
			TransferSyntax: make([]byte, 16),
		},
	}
}

// smb->函数绑定
func (c *SMBClient) MSRPCBind(treeId uint32, fileId []byte, uuid string, version uint32) (err error) {
	header := NewMSRPCHeader()
	header.FragLength = 72
	header.CallId = 1
	header.PacketType = PDUBind
	header.PacketFlags = PDUFlagPending
	bind := MSRPCBindStruct{
		MSRPCHeaderStruct: header,
		MaxXmitFrag:       4280,
		MaxRecvFrag:       4280,
		AssocGroup:        0,
		NumCtxItems:       1,
		CtxItem: CtxEItemStruct{
			NumTransItems: 1,
			AbstractSyntax: SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(uuid),
				Version: version,
			},
			TransferSyntax: SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(NDRSyntax),
				Version: 2,
			},
		},
	}
	req := c.NewWriteRequest(treeId, fileId, bind)
	c.Debug("Sending rpc bind to ["+ms.UUIDMap[uuid]+"]", nil)
	_, err = c.SMBSend(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read rpc response", nil)
	req1 := c.NewReadRequest(treeId, fileId)
	buf, err1 := c.SMBSend(req1)
	if err1 != nil {
		c.Debug("", err1)
		return err1
	}
	smbRes := smb2.NewReadResponse()
	res := NewMSRPCBindAck()
	c.Debug("Unmarshalling rpc bind", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.NumResults < 1 {
		return errors.New("Failed to rpc bind code : [" + ms.UUIDMap[uuid] + "] " + ms.StatusMap[smbRes.Status])
	}
	c.Debug("Completed rpc bind : ["+ms.UUIDMap[uuid]+"]", nil)
	return nil
}

// tcp->函数绑定
func (c *TCPClient) MSRPCBind(uuid string, version uint32) (err error) {
	header := NewMSRPCHeader()
	header.FragLength = 72
	header.CallId = 1
	header.PacketType = PDUBind
	header.PacketFlags = PDUFlagPending
	bind := MSRPCBindStruct{
		MSRPCHeaderStruct: header,
		MaxXmitFrag:       4280,
		MaxRecvFrag:       4280,
		AssocGroup:        0,
		NumCtxItems:       1,
		CtxItem: CtxEItemStruct{
			NumTransItems: 1,
			AbstractSyntax: SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(uuid),
				Version: version,
			},
			TransferSyntax: SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(NDRSyntax),
				Version: 2,
			},
		},
	}
	c.Debug("Sending rpc bind to ["+ms.UUIDMap[uuid]+"]", nil)
	buf, err := c.TCPSend(bind)
	res := NewMSRPCBindAck()
	c.Debug("Unmarshalling rpc bind", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.NumResults < 1 {
		return errors.New("Failed to rpc bind code : [" + ms.UUIDMap[uuid] + "] ")
	}
	c.Debug("Completed rpc bind : ["+ms.UUIDMap[uuid]+"]", nil)
	return err
}
