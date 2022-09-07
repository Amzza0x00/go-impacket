package smb2

import (
	"encoding/hex"
	"errors"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb"
	"go-impacket/pkg/util"
)

// 此文件用于smb2写数据请求
// 将数据写入命名管道、文件

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/e7046961-3318-4350-be2a-a8d69bb59ce8
type SMB2WriteRequestStruct struct {
	smb.SMB2Header
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
	Buffer                 []byte //写入的数据
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/7b80a339-f4d3-4575-8ce2-70a06f24f133
type SMB2WriteResponseStruct struct {
	smb.SMB2Header
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
func (s *Session) NewSMB2WriteRequest(treeId uint32, fileId, buf []byte) SMB2WriteRequestStruct {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_WRITE
	smb2Header.CreditCharge = 1
	smb2Header.Credits = 127
	smb2Header.MessageId = s.messageId
	smb2Header.SessionId = s.sessionId
	smb2Header.TreeId = treeId
	return SMB2WriteRequestStruct{
		SMB2Header:     smb2Header,
		StructureSize:  49,
		FileId:         fileId,
		Channel:        SMB2_CHANNEL_NONE,
		RemainingBytes: 0,
		WriteFlags:     0,
		Buffer:         buf,
	}
}

// 写入请求响应
func NewSMB2WriteResponse() SMB2WriteResponseStruct {
	smb2Header := NewSMB2Header()
	return SMB2WriteResponseStruct{
		SMB2Header: smb2Header,
	}
}

// 需要传入树id
// 一次写入的数据不超过65536
func (s *Session) SMB2WriteRequest(treeId uint32, filepath, filename string, fileId []byte) error {
	s.Debug("Sending Write file request ["+filename+"]", nil)
	// 将文件读入缓冲区
	file, e := util.ReadFile(filepath + filename)
	if e != nil {
		s.Debug("", e)
		return e
	}
	var fileBuff []byte
	// 切分文件大小
	for i := 65536; i < len(file); {
		// 先拿前65536的数据
		fileBuff = file[0 : len(file)-i]
		break
	}
	// 写入第一次切分后的数据
	req := s.NewSMB2WriteRequest(treeId, fileId, fileBuff)
	buf, err := s.send(req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	res := NewSMB2WriteResponse()
	s.Debug("Unmarshalling Write file response ["+filename+"]", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New("Failed to write file to [" + filename + "]: " + ms.StatusMap[res.SMB2Header.Status])
	}
	// 写入第二次切分后的数据
	req = s.NewSMB2WriteRequest(treeId, fileId, file[len(fileBuff):])
	req.FileOffset = uint64(len(fileBuff))
	buf, err = s.send(req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	res = NewSMB2WriteResponse()
	s.Debug("Unmarshalling Write file response ["+filename+"]", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New("Failed to write file to [" + filename + "]: " + ms.StatusMap[res.SMB2Header.Status])
	}
	s.Debug("Completed WriteFile ["+filename+"]", nil)
	return nil
}

// 写入管道数据
func (s *Session) WritePipeRequest(treeId uint32, buffer, fileId []byte) error {
	s.Debug("Sending Write pipe request", nil)
	req := s.NewSMB2WriteRequest(treeId, fileId, buffer)
	buf, err := s.send(req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	res := NewSMB2WriteResponse()
	s.Debug("Unmarshalling Write pipe response", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New("Failed to write pipe to " + ms.StatusMap[res.SMB2Header.Status])
	}
	s.Debug("Completed Write pipe ", nil)
	return nil
}
