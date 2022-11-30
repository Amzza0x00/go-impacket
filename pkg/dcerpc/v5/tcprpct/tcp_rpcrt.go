package tcprpct

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5/dcom"
	"github.com/Amzza0x00/go-impacket/pkg/dcerpc/v5/msrpc"
	"github.com/Amzza0x00/go-impacket/pkg/encoder"
	"github.com/Amzza0x00/go-impacket/pkg/ms"
	"github.com/Amzza0x00/go-impacket/pkg/util"
)

func (c *TCPClient) ServerAlive2Request(callId uint32) (address []string, err error) {
	err = c.MSRPCBind(ms.IID_IObjectExporter, ms.IID_IObjectExporter_VERSION)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	c.Debug("Sending ServerAlive2 request", nil)
	req := dcom.NewServerAlive2Request()
	req.CallId = callId
	buf, err := c.TCPSend(req)
	res := dcom.NewServerAlive2Response()
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

// tcp->函数绑定
func (c *TCPClient) MSRPCBind(uuid string, version uint32) (err error) {
	header := msrpc.NewMSRPCHeader()
	header.FragLength = 72
	header.CallId = 1
	header.PacketType = msrpc.PDUBind
	header.PacketFlags = msrpc.PDUFlagPending
	bind := msrpc.MSRPCBindStruct{
		MSRPCHeaderStruct: header,
		MaxXmitFrag:       4280,
		MaxRecvFrag:       4280,
		AssocGroup:        0,
		NumCtxItems:       1,
		CtxItem: msrpc.CtxEItemStruct{
			NumTransItems: 1,
			AbstractSyntax: msrpc.SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(uuid),
				Version: version,
			},
			TransferSyntax: msrpc.SyntaxIDStruct{
				UUID:    util.PDUUuidFromBytes(msrpc.NDRSyntax),
				Version: 2,
			},
		},
	}
	c.Debug("Sending rpc bind to ["+ms.UUIDMap[uuid]+"]", nil)
	buf, err := c.TCPSend(bind)
	res := msrpc.NewMSRPCBindAck()
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
