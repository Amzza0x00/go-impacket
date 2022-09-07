package smb2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb"
	"strconv"
)

// 此文件用于目录树连接/断开

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/832d2130-22e8-4afb-aafd-b30bb0901798
// 树连接请求结构
type SMB2TreeConnectRequestStruct struct {
	smb.SMB2Header
	StructureSize uint16
	Reserved      uint16 //2字节，smb3.x使用，其他忽略
	PathOffset    uint16 `smb:"offset:Path"`
	PathLength    uint16 `smb:"len:Path"`
	Path          []byte
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/dd34e26c-a75e-47fa-aab2-6efc27502e96
// 树连接响应结构
type SMB2TreeConnectResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	ShareType     uint8 //1字节，访问共享类型
	Reserved      uint8 //1字节
	ShareFlags    uint32
	Capabilities  uint32
	MaximalAccess uint32
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/8a622ecb-ffee-41b9-b4c4-83ff2d3aba1b
// 断开树连接请求结构
type SMB2TreeDisconnectRequestStruct struct {
	smb.SMB2Header
	StructureSize uint16 //2字节，客户端必须设为4,表示请求大小
	Reserved      uint16
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/aeac92de-8db3-48f8-a8b7-bfee28b9fd9e
// 断开树连接响应结构
type SMB2TreeDisconnectResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	Reserved      uint16
}

func (s *Session) NewSMB2TreeConnectRequest(name string) (SMB2TreeConnectRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_TREE_CONNECT
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = s.messageId
	smb2Header.SessionId = s.sessionId

	//格式 \\172.20.10.5:445\IPC$
	path := fmt.Sprintf("\\\\%s:%s\\%s", s.options.Host, strconv.Itoa(s.options.Port), name)
	return SMB2TreeConnectRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 9,
		Reserved:      0,
		PathOffset:    0,
		PathLength:    0,
		Path:          encoder.ToUnicode(path),
	}, nil
}

func NewSMB2TreeConnectResponse() SMB2TreeConnectResponseStruct {
	smb2Header := NewSMB2Header()
	return SMB2TreeConnectResponseStruct{
		SMB2Header: smb2Header,
	}
}

func (s *Session) NewSMB2TreeDisconnectRequest(treeId uint32) (SMB2TreeDisconnectRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_TREE_DISCONNECT
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = s.messageId
	smb2Header.SessionId = s.sessionId
	smb2Header.TreeId = treeId

	return SMB2TreeDisconnectRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 4,
		Reserved:      0,
	}, nil
}

func NewSMB2TreeDisconnectResponse() SMB2TreeDisconnectResponseStruct {
	smb2Header := NewSMB2Header()
	return SMB2TreeDisconnectResponseStruct{
		SMB2Header: smb2Header,
	}
}

// 树连接
func (s *Session) SMB2TreeConnect(name string) (treeId uint32, err error) {
	s.Debug("Sending TreeConnect request ["+name+"]", nil)
	req, err := s.NewSMB2TreeConnectRequest(name)
	if err != nil {
		s.Debug("", err)
		return 0, err
	}
	buf, err := s.send(req)
	if err != nil {
		s.Debug("", err)
		return 0, err
	}
	res := NewSMB2TreeConnectResponse()
	s.Debug("Unmarshalling TreeConnect response ["+name+"]", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
		//return err
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return 0, errors.New("Failed to connect to [" + name + "]: " + ms.StatusMap[res.SMB2Header.Status])
	}
	treeID := res.SMB2Header.TreeId
	s.trees[name] = treeID
	s.Debug("Completed TreeConnect ["+name+"]", nil)
	return treeID, nil
}

// 断开树连接
func (s *Session) SMB2TreeDisconnect(name string) error {
	var (
		treeid    uint32
		pathFound bool
	)
	for k, v := range s.trees {
		if k == name {
			treeid = v
			pathFound = true
			break
		}
	}
	if !pathFound {
		err := errors.New("Unable to find tree path for disconnect")
		s.Debug("", err)
		return err
	}
	s.Debug("Sending TreeDisconnect request ["+name+"]", nil)
	req, err := s.NewSMB2TreeDisconnectRequest(treeid)
	if err != nil {
		s.Debug("", err)
		return err
	}
	buf, err := s.send(req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	s.Debug("Unmarshalling TreeDisconnect response for ["+name+"]", nil)
	res := NewSMB2TreeDisconnectResponse()
	if err := encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if res.SMB2Header.Status != ms.STATUS_MORE_PROCESSING_REQUIRED {
		return errors.New("Failed to connect to tree: " + ms.StatusMap[res.SMB2Header.Status])
	}
	delete(s.trees, name)
	s.Debug("TreeDisconnect completed ["+name+"]", nil)
	return nil
}
