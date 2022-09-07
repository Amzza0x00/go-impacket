package smb2

import (
	"encoding/hex"
	"errors"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb"
	"go-impacket/pkg/util"
)

// 此文件提供ms-rpc封装
// DCE/RPC RPC over SMB 协议实现
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-srvs/accf23b0-0f57-441c-9185-43041f1b0ee9
// https://pubs.opengroup.org/onlinepubs/9629399/

// RPC over SMB 标准头
type PDUHeader struct {
	smb.SMB2Header
	StructureSize          uint16
	DataOffset             uint16 `smb:"offset:Buffer"`
	WriteLength            uint32 `smb:"len:Buffer"`
	FileOffset             []byte `smb:"fixed:8"`
	FileId                 []byte `smb:"fixed:16"` //16字节，服务端返回句柄
	Channel                uint32
	RemainingBytes         uint32
	WriteChannelInfoOffset uint16
	WriteChannelInfoLength uint16
	WriteFlags             uint32
	Buffer                 interface{} //写入的数据
}

// DCE/RPC 标准头
type PDUHeaderStruct struct {
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32 //4字节，小端排序，0x10
	FragLength         uint16 //2字节，整个结构的长度
	AuthLength         uint16
	CallId             uint32
	Buffer             interface{}
}

// 函数绑定结构
type PDUBindStruct struct {
	//PDUHeader
	MaxXmitFrag uint16 //4字节，发送大小协商
	MaxRecvFrag uint16 //4字节，接收大小协商
	AssocGroup  uint32
	NumCtxItems uint8
	Reserved    uint8
	Reserved2   uint16
	CtxItem     PDUCtxEItem
}

// PDU CtxItem结构
type PDUCtxEItem struct {
	ContextId      uint16
	NumTransItems  uint8
	Reserved       uint8
	AbstractSyntax PDUSyntaxID
	TransferSyntax PDUSyntaxID
}

type PDUSyntaxID struct {
	UUID    []byte `smb:"fixed:16"`
	Version uint32
}

// PDU CtxItem响应结构
type PDUCtxEItemResponseStruct struct {
	AckResult      uint16
	AckReason      uint16
	TransferSyntax []byte `smb:"fixed:16"` //16字节
	SyntaxVer      uint32
}

type PDUBindAckStruct struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	MaxXmitFrag        uint16
	MaxRecvFrag        uint16
	AssocGroup         uint32
	ScndryAddrlen      uint16
	ScndryAddr         []byte `smb:"count:ScndryAddrlen"` //取决管道的长度
	NumResults         uint8
	CtxItem            PDUCtxEItemResponseStruct
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

func NewPDUHeader() PDUHeader {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_WRITE
	smb2Header.CreditCharge = 1
	return PDUHeader{
		SMB2Header:    smb2Header,
		StructureSize: 49,
		//WriteLength:            72,
		FileOffset:             make([]byte, 8),
		Channel:                SMB2_CHANNEL_NONE,
		RemainingBytes:         0,
		WriteChannelInfoOffset: 0,
		WriteChannelInfoLength: 0,
		WriteFlags:             0,
		//Buffer: PDUHeaderStruct{
		//	Version:            5,
		//	VersionMinor:       0,
		//	DataRepresentation: 16,
		//	AuthLength:         0,
		//	CallId:             1,
		//},
	}
}

// 函数绑定请求
func (s *Session) NewSMB2PDUBind(treeId uint32, fileId []byte, uuid string, version uint32) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	pduHeader.Buffer = PDUHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDUBind,
		PacketFlags:        PDUFlagPending,
		DataRepresentation: 16,
		FragLength:         72,
		AuthLength:         0,
		CallId:             1,
		Buffer: PDUBindStruct{
			MaxXmitFrag: 4280,
			MaxRecvFrag: 4280,
			AssocGroup:  0,
			NumCtxItems: 1,
			CtxItem: PDUCtxEItem{
				ContextId:     0,
				NumTransItems: 1,
				Reserved:      0,
				AbstractSyntax: PDUSyntaxID{
					UUID:    util.PDUUuidFromBytes(uuid),
					Version: version,
				},
				TransferSyntax: PDUSyntaxID{
					UUID:    util.PDUUuidFromBytes(NDRSyntax),
					Version: 2,
				},
			}},
	}
	//pduHeader.Buffer.PacketType = PDUBind
	//pduHeader.Buffer.PacketFlags = PDUFlagPending
	//// 设置函数绑定
	//pduHeader.Buffer.Buffer = PDUBindStruct{
	//	MaxXmitFrag: 4280,
	//	MaxRecvFrag: 4280,
	//	AssocGroup:  0,
	//	NumCtxItems: 1,
	//	CtxItem: PDUCtxEItem{
	//		ContextId:     0,
	//		NumTransItems: 1,
	//		Reserved:      0,
	//		AbstractSyntax: PDUSyntaxID{
	//			UUID:    util.PDUUuidFromBytes(uuid),
	//			Version: 3,
	//		},
	//		TransferSyntax: PDUSyntaxID{
	//			UUID:    util.PDUUuidFromBytes(NDRSyntax),
	//			Version: 2,
	//		},
	//	}}
	//// 正常情况下固定为72
	//pduHeader.Buffer.FragLength = 72
	return pduHeader
	//return pduHeader{
	//	Buffer: PDUHeaderStruct{
	//		PacketType:  PDUBind,
	//		PacketFlags: PDUFlagPending,
	//		FragLength:  72,
	//		Buffer: PDUBindStruct{
	//			MaxXmitFrag: 4280,
	//			MaxRecvFrag: 4280,
	//			AssocGroup:  0,
	//			NumCtxItems: 1,
	//			CtxItem: PDUCtxEItem{
	//				ContextId:     0,
	//				NumTransItems: 1,
	//				Reserved:      0,
	//				AbstractSyntax: PDUSyntaxID{
	//					UUID:    util.PDUUuidFromBytes(uuid),
	//					Version: 3,
	//				},
	//				TransferSyntax: PDUSyntaxID{
	//					UUID:    util.PDUUuidFromBytes(NDRSyntax),
	//					Version: 2,
	//				},
	//			},
	//		},
	//	},
	//}
	//smb2Header := NewSMB2Header()
	//smb2Header.Command = smb.SMB2_WRITE
	//smb2Header.CreditCharge = 1
	//smb2Header.messageId = s.messageId
	//smb2Header.sessionId = s.sessionId
	//smb2Header.TreeId = treeId
	//return PDUHeader{
	//	SMB2Header:             smb2Header,
	//	StructureSize:          49,
	//	WriteLength:            72,
	//	FileOffset:             make([]byte, 8),
	//	FileId:                 fileId,
	//	Channel:                SMB2_CHANNEL_NONE,
	//	RemainingBytes:         0,
	//	WriteChannelInfoOffset: 0,
	//	WriteChannelInfoLength: 0,
	//	WriteFlags:             0,
	//	Buffer: PDUHeaderStruct{
	//		Version:            5,
	//		VersionMinor:       0,
	//		PacketType:         PDUBind,
	//		PacketFlags:        PDUFlagPending,
	//		DataRepresentation: 16,
	//		FragLength:         72,
	//		AuthLength:         0,
	//		CallId:             1,
	//		Buffer: PDUBindStruct{
	//			MaxXmitFrag: 4280,
	//			MaxRecvFrag: 4280,
	//			AssocGroup:  0,
	//			NumCtxItems: 1,
	//			CtxItem: PDUCtxEItem{
	//				ContextId:     0,
	//				NumTransItems: 1,
	//				Reserved:      0,
	//				AbstractSyntax: PDUSyntaxID{
	//					UUID:    util.PDUUuidFromBytes("4b324fc8-1670-01d3-1278-5a47bf6ee188"), // svcctl uuid
	//					Version: 3,
	//				},
	//				TransferSyntax: PDUSyntaxID{
	//					UUID:    util.PDUUuidFromBytes(NDRSyntax),
	//					Version: 2,
	//				},
	//			},
	//		},
	//	},
	//}
}

// 函数绑定响应
func NewSMB2PDUBindAck() PDUBindAckStruct {
	smb2Header := NewSMB2Header()
	return PDUBindAckStruct{
		SMB2ReadResponseStruct: SMB2ReadResponseStruct{
			SMB2Header: smb2Header,
		},
		CtxItem: PDUCtxEItemResponseStruct{
			TransferSyntax: make([]byte, 16),
		},
	}
}

func (s *Session) SMB2PDUBind(treeId uint32, fileId []byte, uuid string, version uint32) error {
	s.Debug("Sending rpc bind to ["+ms.UUIDMap[uuid]+"]", nil)
	req := s.NewSMB2PDUBind(treeId, fileId, uuid, version)
	_, err := s.send(req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	s.Debug("Read rpc response", nil)
	req1 := s.NewSMB2ReadRequest(treeId, fileId)
	buf, err1 := s.send(req1)
	if err1 != nil {
		s.Debug("", err1)
		return err1
	}
	res := NewSMB2PDUBindAck()
	s.Debug("Unmarshalling rpc bind", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New("Failed to rpc bind to [" + ms.UUIDMap[uuid] + "] " + ms.StatusMap[res.SMB2Header.Status])
	}
	s.Debug("Completed rpc bind to ["+ms.UUIDMap[uuid]+"]", nil)
	return nil
}
