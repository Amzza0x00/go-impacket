package smb2

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/gss"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/ntlm"
	"go-impacket/pkg/smb"
	"io"
	"log"
	"net"
	"runtime/debug"
)

// 此文件提供smb连接方法

// SMB会话结构
type Session struct {
	IsSigningRequired bool
	IsAuthenticated   bool
	debug             bool
	securityMode      uint16
	messageId         uint64
	sessionId         uint64
	conn              net.Conn
	dialect           uint16
	options           Options
	trees             map[string]uint32
}

// SMB连接参数
type Options struct {
	Host        string
	Port        int
	Workstation string
	Domain      string
	User        string
	Password    string
	Hash        string
}

func (s *Session) Debug(msg string, err error) {
	if s.debug {
		log.Println("[ DEBUG ] ", msg)
		if err != nil {
			debug.PrintStack()
		}
	}
}

func NewSMB2Header() smb.SMB2Header {
	return smb.SMB2Header{
		ProtocolId:    []byte(smb.ProtocolSMB2),
		StructureSize: 64,
		CreditCharge:  0,
		Status:        0,
		Command:       0,
		Credits:       0,
		Flags:         0,
		NextCommand:   0,
		MessageId:     0,
		Reserved:      0,
		TreeId:        0,
		SessionId:     0,
		Signature:     make([]byte, 16),
	}
}

// 协商版本请求初始化
func (s *Session) NewSMB2NegotiateRequest() smb.SMB2NegotiateRequestStruct {
	// 初始化
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_NEGOTIATE
	smb2Header.MessageId = s.messageId
	smb2Header.CreditCharge = 1
	return smb.SMB2NegotiateRequestStruct{
		SMB2Header:      smb2Header,
		StructureSize:   36,
		DialectCount:    1,
		SecurityMode:    smb.SecurityModeSigningEnabled, // 必须开启签名
		Reserved:        0,
		Capabilities:    0,
		ClientGuid:      make([]byte, 16),
		ClientStartTime: 0,
		Dialects: []uint16{
			uint16(smb.SMB2_1_Dialect),
		},
	}
}

// 协商版本响应初始化
func NewSMB2NegotiateResponse() smb.SMB2NegotiateResponseStruct {
	smb2Header := NewSMB2Header()
	return smb.SMB2NegotiateResponseStruct{
		SMB2Header:           smb2Header,
		StructureSize:        0,
		SecurityMode:         0,
		DialectRevision:      0,
		Reserved:             0,
		ServerGuid:           make([]byte, 16),
		Capabilities:         0,
		MaxTransactSize:      0,
		MaxReadSize:          0,
		MaxWriteSize:         0,
		SystemTime:           0,
		ServerStartTime:      0,
		SecurityBufferOffset: 0,
		SecurityBufferLength: 0,
		Reserved2:            0,
		SecurityBlob:         &gss.NegTokenInit{},
	}
}

// 质询请求初始化
func (s *Session) NewSMB2SessionSetupRequest() (smb.SMB2SessionSetupRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_SESSION_SETUP
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = s.messageId
	smb2Header.SessionId = s.sessionId

	ntlmsspneg := ntlm.NewNegotiate(s.options.Domain, s.options.Workstation)
	data, err := encoder.Marshal(ntlmsspneg)
	if err != nil {
		return smb.SMB2SessionSetupRequestStruct{}, err
	}

	if s.sessionId != 0 {
		return smb.SMB2SessionSetupRequestStruct{}, errors.New("Bad session ID for session setup 1 message")
	}

	// Initial session setup request
	init, err := gss.NewNegTokenInit()
	if err != nil {
		return smb.SMB2SessionSetupRequestStruct{}, err
	}
	init.Data.MechToken = data

	return smb.SMB2SessionSetupRequestStruct{
		SMB2Header:           smb2Header,
		StructureSize:        25,
		Flags:                0x00,
		SecurityMode:         byte(smb.SecurityModeSigningEnabled),
		Capabilities:         0,
		Channel:              0,
		SecurityBufferOffset: 88,
		SecurityBufferLength: 0,
		PreviousSessionID:    0,
		SecurityBlob:         &init,
	}, nil
}

// 质询响应初始化
func NewSMB2SessionSetupResponse() (smb.SMB2SessionSetupResponseStruct, error) {
	smb2Header := NewSMB2Header()
	resp, err := gss.NewNegTokenResp()
	if err != nil {
		return smb.SMB2SessionSetupResponseStruct{}, err
	}
	ret := smb.SMB2SessionSetupResponseStruct{
		SMB2Header:   smb2Header,
		SecurityBlob: &resp,
	}
	return ret, nil
}

// 认证请求初始化
func (s *Session) NewSMB2SessionSetup2Request() (smb.SMB2SessionSetup2RequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_SESSION_SETUP
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = s.messageId
	smb2Header.SessionId = s.sessionId

	ntlmsspneg := ntlm.NewNegotiate(s.options.Domain, s.options.Workstation)
	data, err := encoder.Marshal(ntlmsspneg)
	if err != nil {
		return smb.SMB2SessionSetup2RequestStruct{}, err
	}

	if s.sessionId == 0 {
		return smb.SMB2SessionSetup2RequestStruct{}, errors.New("Bad session ID for session setup 2 message")
	}

	// Session setup request #2
	resp, err := gss.NewNegTokenResp()
	if err != nil {
		return smb.SMB2SessionSetup2RequestStruct{}, err
	}
	resp.ResponseToken = data

	return smb.SMB2SessionSetup2RequestStruct{
		SMB2Header:           smb2Header,
		StructureSize:        25,
		Flags:                0x00,
		SecurityMode:         byte(smb.SecurityModeSigningEnabled),
		Capabilities:         0,
		Channel:              0,
		SecurityBufferOffset: 88,
		SecurityBufferLength: 0,
		PreviousSessionID:    0,
		SecurityBlob:         &resp,
	}, nil
}

func (s *Session) SMB2NegotiateProtocol() error {
	// 第一步 发送协商请求
	s.Debug("Sending Negotiate request", nil)
	negReq := s.NewSMB2NegotiateRequest()
	buf, err := s.send(negReq)
	if err != nil {
		s.Debug("", err)
		return err
	}
	negRes := NewSMB2NegotiateResponse()
	if err = encoder.Unmarshal(buf, &negRes); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if negRes.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New(fmt.Sprintf("NT Status Error: %d\n", negRes.SMB2Header.Status))
	}
	// Check SPNEGO security blob
	//spnegoOID, err := encoder.ObjectIDStrToInt(encoder.SpnegoOid)
	//if err != nil {
	//	s.Debug(err)
	//	return err
	//}
	//oid := negRes.SecurityBlob.OID
	//fmt.Println(oid)
	// 检查是否存在ntlmssp
	hasNTLMSSP := false
	ntlmsspOID, err := gss.ObjectIDStrToInt(ntlm.NTLMSSPMECHTYPEOID)
	if err != nil {
		return err
	}
	for _, mechType := range negRes.SecurityBlob.Data.MechTypes {
		if mechType.Equal(ntlmsspOID) {
			hasNTLMSSP = true
			break
		}
	}
	if !hasNTLMSSP {
		return errors.New("Server does not support NTLMSSP")
	}
	// 设置会话安全模式
	s.securityMode = negRes.SecurityMode
	// 设置会话协议
	s.dialect = negRes.DialectRevision
	// 签名开启/关闭
	mode := s.securityMode
	if mode&smb.SecurityModeSigningEnabled > 0 {
		if mode&smb.SecurityModeSigningRequired > 0 {
			s.IsSigningRequired = true
		} else {
			s.IsSigningRequired = false
		}
	} else {
		s.IsSigningRequired = false
	}
	// 第二步 发送质询
	s.Debug("Sending SessionSetup1 request", nil)
	ssreq, err := s.NewSMB2SessionSetupRequest()
	if err != nil {
		s.Debug("", err)
		return err
	}
	ssres, err := NewSMB2SessionSetupResponse()
	if err != nil {
		s.Debug("", err)
		return err
	}
	buf, err = encoder.Marshal(ssreq)
	if err != nil {
		s.Debug("", err)
		return err
	}

	buf, err = s.send(ssreq)
	if err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}

	s.Debug("Unmarshalling SessionSetup1 response", nil)
	if err = encoder.Unmarshal(buf, &ssres); err != nil {
		s.Debug("", err)
		return err
	}

	challenge := ntlm.NewChallenge()
	resp := ssres.SecurityBlob
	if err = encoder.Unmarshal(resp.ResponseToken, &challenge); err != nil {
		s.Debug("", err)
		return err
	}

	if ssres.SMB2Header.Status != ms.STATUS_MORE_PROCESSING_REQUIRED {
		status, _ := ms.StatusMap[negRes.SMB2Header.Status]
		return errors.New(fmt.Sprintf("NT Status Error: %s\n", status))
	}
	s.sessionId = ssres.SMB2Header.SessionId

	s.Debug("Sending SessionSetup2 request", nil)
	// 第三步 认证
	ss2req, err := s.NewSMB2SessionSetup2Request()
	if err != nil {
		s.Debug("", err)
		return err
	}

	var auth ntlm.NTLMv2Authentication
	if s.options.Hash != "" {
		// Hash present, use it for auth
		s.Debug("Performing hash-based authentication", nil)
		auth = ntlm.NewAuthenticateHash(s.options.Domain, s.options.User, s.options.Workstation, s.options.Hash, challenge)
	} else {
		// No hash, use password
		s.Debug("Performing password-based authentication", nil)
		auth = ntlm.NewAuthenticatePass(s.options.Domain, s.options.User, s.options.Workstation, s.options.Password, challenge)
	}

	responseToken, err := encoder.Marshal(auth)
	if err != nil {
		s.Debug("", err)
		return err
	}
	resp2 := ss2req.SecurityBlob
	resp2.ResponseToken = responseToken
	ss2req.SecurityBlob = resp2
	ss2req.SMB2Header.Credits = 127
	buf, err = encoder.Marshal(ss2req)
	if err != nil {
		s.Debug("", err)
		return err
	}

	buf, err = s.send(ss2req)
	if err != nil {
		s.Debug("", err)
		return err
	}
	s.Debug("Unmarshalling SessionSetup2 response", nil)
	var authResp smb.SMB2Header
	if err = encoder.Unmarshal(buf, &authResp); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if authResp.Status != ms.STATUS_SUCCESS {
		// authResp.Status 十进制表示
		status, _ := ms.StatusMap[authResp.Status]
		return errors.New(fmt.Sprintf("NT Status Error: %s\n", status))
	}
	s.IsAuthenticated = true

	s.Debug("Completed NegotiateProtocol and SessionSetup", nil)
	return nil
}

// SMB2连接封装
func NewSession(opt Options, debug bool) (s *Session, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
	if err != nil {
		return
	}

	s = &Session{
		IsSigningRequired: false,
		IsAuthenticated:   false,
		debug:             debug,
		securityMode:      0,
		messageId:         0,
		sessionId:         0,
		dialect:           0,
		conn:              conn,
		options:           opt,
		trees:             make(map[string]uint32),
	}

	err = s.SMB2NegotiateProtocol()
	if err != nil {
		return
	}

	return s, nil
}

func (s *Session) send(req interface{}) (res []byte, err error) {
	buf, err := encoder.Marshal(req)
	if err != nil {
		s.Debug("", err)
		return nil, err
	}

	b := new(bytes.Buffer)
	if err = binary.Write(b, binary.BigEndian, uint32(len(buf))); err != nil {
		s.Debug("", err)
		return
	}
	s.Debug("Raw:\n"+hex.Dump(append(b.Bytes(), buf...)), nil)
	rw := bufio.NewReadWriter(bufio.NewReader(s.conn), bufio.NewWriter(s.conn))
	if _, err = rw.Write(append(b.Bytes(), buf...)); err != nil {
		s.Debug("", err)
		return
	}
	rw.Flush()

	var size uint32
	if err = binary.Read(rw, binary.BigEndian, &size); err != nil {
		s.Debug("", err)
		return
	}
	if size > 0x00FFFFFF {
		return nil, errors.New("Invalid NetBIOS Session message")
	}

	data := make([]byte, size)
	l, err := io.ReadFull(rw, data)
	if err != nil {
		s.Debug("", err)
		return nil, err
	}
	if uint32(l) != size {
		return nil, errors.New("Message size invalid")
	}

	//protID := data[0:4]
	//switch string(protID) {
	//default:
	//	return nil, errors.New("Protocol Not Implemented")
	//case ProtocolSMB:
	//}

	s.messageId++
	return data, nil
}

func (s *Session) Close() {
	s.Debug("Closing session", nil)
	for k, _ := range s.trees {
		s.SMB2TreeDisconnect(k)
	}
	s.conn.Close()
	s.Debug("Session close completed", nil)
}
