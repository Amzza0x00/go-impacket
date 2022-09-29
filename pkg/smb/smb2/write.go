package smb2

import (
	"encoding/hex"
	"errors"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb"
	"os"
)

// 此文件用于smb2写数据请求
// 将数据写入命名管道、文件

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/e7046961-3318-4350-be2a-a8d69bb59ce8
type WriteRequestStruct struct {
	smb.SMB2PacketStruct
	StructureSize          uint16
	DataOffset             uint16 `smb:"offset:Buffer"`
	WriteLength            uint32 `smb:"len:Buffer"`
	FileOffset             uint64
	FileId                 []byte `smb:"fixed:16"` //16字节，服务端返回句柄
	Channel                uint32
	RemainingBytes         uint32
	WriteChannelInfoOffset uint16
	WriteChannelInfoLength uint16
	WriteFlags             uint32
	Buffer                 interface{} //写入的数据
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/7b80a339-f4d3-4575-8ce2-70a06f24f133
type WriteResponseStruct struct {
	smb.SMB2PacketStruct
	StructureSize          uint16
	Reserved               uint16
	WriteCount             uint32
	WriteRemaining         uint32
	WriteChannelInfoOffset uint16
	WriteChannelInfoLength uint16
}

// Channel属性
const (
	SMB2_CHANNEL_NONE               = 0x00000000
	SMB2_CHANNEL_RDMA_V1            = 0x00000001
	SMB2_CHANNEL_RDMA_V1_INVALIDATE = 0x00000002
	SMB2_CHANNEL_RDMA_TRANSFORM     = 0x00000003
)

// 写入请求
func (c *Client) NewWriteRequest(treeId uint32, fileId []byte, buf interface{}) WriteRequestStruct {
	smb2Header := NewSMB2Packet()
	smb2Header.Command = smb.SMB2_WRITE
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.TreeId = treeId
	return WriteRequestStruct{
		SMB2PacketStruct: smb2Header,
		StructureSize:    49,
		FileId:           fileId,
		Channel:          SMB2_CHANNEL_NONE,
		RemainingBytes:   0,
		WriteFlags:       0,
		Buffer:           buf,
	}
}

// 写入请求响应
func NewWriteResponse() WriteResponseStruct {
	smb2Header := NewSMB2Packet()
	return WriteResponseStruct{
		SMB2PacketStruct: smb2Header,
	}
}

// 需要传入树id
func (c *Client) WriteRequest(treeId uint32, filepath, filename string, fileId []byte) (err error) {
	c.Debug("Sending Write file request ["+filename+"]", nil)
	// 将文件读入缓冲区
	file, err := os.Open(filepath + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// 一次传入1kb数据
	fileBuf := make([]byte, 10240)
	fileOffset := 0
	i := 0
Loop:
	for {
		switch nr, _ := file.Read(fileBuf[:]); true {
		case nr < 0:
			return errors.New("Failed read file to [" + filepath + filename + "]")
		case nr == 0: // EOF
			break Loop
		case nr > 0:
			req := c.NewWriteRequest(treeId, fileId, fileBuf)
			if i == 0 {
				req.FileOffset = 0
			} else {
				req.FileOffset = uint64(fileOffset)
			}
			fileOffset += nr
			i++
			buf, err := c.Send(req)
			if err != nil {
				c.Debug("", err)
				return err
			}
			res := NewWriteResponse()
			c.Debug("Unmarshalling Write file response ["+filename+"]", nil)
			if err = encoder.Unmarshal(buf, &res); err != nil {
				c.Debug("Raw:\n"+hex.Dump(buf), err)
			}
			if res.SMB2PacketStruct.Status != ms.STATUS_SUCCESS {
				return errors.New("Failed to write file to [" + filename + "]: " + ms.StatusMap[res.SMB2PacketStruct.Status])
			}
		}
	}
	c.Debug("Completed WriteFile ["+filename+"]", nil)
	return nil
}

// 写入管道数据
func (c *Client) WritePipeRequest(treeId uint32, buffer, fileId []byte) error {
	c.Debug("Sending Write pipe request", nil)
	req := c.NewWriteRequest(treeId, fileId, buffer)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	res := NewWriteResponse()
	c.Debug("Unmarshalling Write pipe response", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2PacketStruct.Status != ms.STATUS_SUCCESS {
		return errors.New("Failed to write pipe to " + ms.StatusMap[res.SMB2PacketStruct.Status])
	}
	c.Debug("Completed Write pipe ", nil)
	return nil
}
