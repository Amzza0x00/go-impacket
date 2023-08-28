package v5

import (
	"encoding/hex"
	"github.com/Amzza0x00/go-impacket/pkg/encoder"
	"github.com/Amzza0x00/go-impacket/pkg/ms"
	"github.com/Amzza0x00/go-impacket/pkg/util"
)

// 此文件提供epmapper rpc接口
// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-rpce/86fc67d3-f44c-4a14-afeb-1e46048841c5

// 绑定epmapper接口
func (c *TCPClient) RpcBindEpmapper(callId uint32) (err error) {
	ctxs := []CtxItemStruct{{
		NumTransItems: 1,
		AbstractSyntax: SyntaxIDStruct{
			UUID:    util.PDUUuidFromBytes(ms.EPMv4_UUID),
			Version: ms.EPMv4_VERSION,
		},
		TransferSyntax: SyntaxIDStruct{
			UUID:    util.PDUUuidFromBytes(ms.NDR_UUID),
			Version: ms.NDR_VERSION,
		}}}
	_, err = c.MSRPCBind(callId, ctxs)
	if err != nil {
		c.Debug("", err)
		return err
	}
	return nil
}

// lookup request请求结构
type EPMLookupRequestStruct struct {
	MSRPCHeaderStruct
	AllocHint            uint32 `smb:"len:EndpointMapperLookup"` // Endpoint Mapper lookup长度
	ContextId            uint16
	Opnum                uint16
	EndpointMapperLookup endpointMapperLookup
}

type endpointMapperLookup struct {
	InquiryType          uint32
	NullPointerObject    uint32
	NullPointerInterface uint32
	VersionOption        uint32
	EntryHandle          []byte `smb:"fixed:20"`
	MaxEntries           uint32
}

func NewEPMLookupRequest() EPMLookupRequestStruct {
	header := NewMSRPCHeader()
	header.PacketType = PDURequest
	header.PacketFlags = PDUFault
	header.FragLength = 24
	return EPMLookupRequestStruct{
		MSRPCHeaderStruct: header,
		Opnum:             2,
		EndpointMapperLookup: endpointMapperLookup{
			InquiryType:          0,
			NullPointerObject:    0,
			NullPointerInterface: 0,
			VersionOption:        1,
			EntryHandle:          make([]byte, 20),
			MaxEntries:           500,
		},
	}
}

type EPMLookupResponseStruct struct {
	MSRPCHeaderStruct
	AllocHint   uint32
	ContextId   uint16
	CancelCount uint8
	Reserved    uint8
	// epm结构
	EntryHandle        []byte `smb:"fixed:20"`
	NumEntries         uint32
	EntriesMaxCount    uint32
	EntriesOffset      uint32
	EntriesActualCount uint32
	// 块结构
	Entry     []EntryStruct `smb:"value:EntriesActualCount"`
	Reserved2 []byte        `smb:"fixed:3"`
	// 块的具体数据
	EntryTowerPointer []EntryTowerPointerStruct `smb:"value:EntriesActualCount"`
	ReturnCode        uint32
}

type EntryStruct struct {
	Object           []byte `smb:"fixed:16"`
	ReferentID       []byte `smb:"fixed:4"`
	AnnotationOffset uint32
	AnnotationLength uint32
	Annotation       []byte `smb:"dynamic:AnnotationLength:4"` // 根据AnnotationLength动态设置长度，最小长度4字节，并且如果长度不能整除4就填充00直到能整除
}

type EntryTowerPointerStruct struct {
	Length         uint32 // 数据总长度
	Length1        uint32
	NumberOfFloors uint16 // 数量
	Buffer         []byte `smb:"value:Length"` // 可变的
}

func NewEPMLookupResponse() EPMLookupResponseStruct {
	return EPMLookupResponseStruct{}
}

func (c *TCPClient) EPMLookupRequest(callId uint32) (res EPMLookupResponseStruct, err error) {
	c.Debug("Sending EPM Lookup request", nil)
	req := NewEPMLookupRequest()
	fragLength := util.SizeOfStruct(req.EndpointMapperLookup)
	req.FragLength += uint16(fragLength)
	req.CallId = callId
	buf, err := c.TCPSend(req)
	// 解析响应内容
	res = NewEPMLookupResponse()
	c.Debug("Unmarshalling EPMLookup response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	return res, nil
}
