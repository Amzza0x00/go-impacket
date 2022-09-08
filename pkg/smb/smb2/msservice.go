package smb2

import (
	"encoding/hex"
	"errors"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/util"
)

// 此文件提供访问windows服务管理封装

// DCE/RPC 扩展头
// 调用win ms service control api
type PDUExtHeaderStruct struct {
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32 //4字节，小端排序，0x10
	FragLength         uint16 //2字节，整个结构的长度
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32 `smb:"len:Buffer"` //Buffer的长度
	ContextId          uint16
	OpNum              uint16
	Buffer             interface{}
}

// ms service control
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/dc84adb3-d51d-48eb-820d-ba1c6ca5faf2
type OpenSCManagerWStruct struct {
	MachineName machineName
	Database    database
	AccessMask  uint32
}

type machineName struct {
	ReferentId  uint32 `smb:"offset:MachineName"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32 //机器名的长度
	MachineName []byte //任意长度,unicode编码
	Reserved    uint16
}

type database struct {
	ReferentId  uint32 `smb:"offset:Database"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32 //机器名的长度
	Database    []byte //任意长度,unicode编码
	Reserved    uint16
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/0d7a7011-9f41-470d-ad52-8535b47ac282
// 安全描述符
const (
	SERVICE_ALL_ACCESS        = 0x000F01FF
	SC_MANAGER_CREATE_SERVICE = 0x00000002
	SC_MANAGER_CONNECT        = 0x00000001
)

// OpenSCManagerW响应结构
type OpenSCManagerWResponse struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

// opnum
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/0d7a7011-9f41-470d-ad52-8535b47ac282
const (
	RCloseServiceHandle         = 0
	RControlService             = 1
	RDeleteService              = 2
	RLockServiceDatabase        = 3
	RQueryServiceObjectSecurity = 4
	RSetServiceObjectSecurity   = 5
	RQueryServiceStatus         = 6
	RSetServiceStatus           = 7
	RUnlockServiceDatabase      = 8
	RNotifyBootConfigStatus     = 9
	RChangeServiceConfigW       = 11
	RCreateServiceW             = 12
	REnumDependentServicesW     = 13
	REnumServicesStatusW        = 14
	ROpenSCManagerW             = 15
	ROpenServiceW               = 16
	RQueryServiceConfigW        = 17
	RQueryServiceLockStatusW    = 18
	RStartServiceW              = 19
	RGetServiceDisplayNameW     = 20
	RGetServiceKeyNameW         = 21
	RChangeServiceConfigA       = 23
	RCreateServiceA             = 24
	REnumDependentServicesA     = 25
	REnumServicesStatusA        = 26
	ROpenSCManagerA             = 27
	ROpenServiceA               = 28
	RQueryServiceConfigA        = 29
	RQueryServiceLockStatusA    = 30
	RStartServiceA              = 31
	RGetServiceDisplayNameA     = 32
	RGetServiceKeyNameA         = 33
	REnumServiceGroupW          = 35
	RChangeServiceConfig2A      = 36
	RChangeServiceConfig2W      = 37
	RQueryServiceConfig2A       = 38
	RQueryServiceConfig2W       = 39
	RQueryServiceStatusEx       = 40
	REnumServicesStatusExA      = 41
	REnumServicesStatusExW      = 42
	RCreateServiceWOW64A        = 44
	RCreateServiceWOW64W        = 45
	RNotifyServiceStatusChange  = 47
	RGetNotifyResults           = 48
	RCloseNotifyHandle          = 49
	RControlServiceExA          = 50
	RControlServiceExW          = 51
	RQueryServiceConfigEx       = 56
	RCreateWowService           = 60
	ROpenSCManager2             = 64
)

// OpenSCManagerW请求
//DWORD ROpenSCManagerW(
//	[in, string, unique, range(0, SC_MAX_COMPUTER_NAME_LENGTH)] SVCCTL_HANDLEW lpMachineName,
//	[in, string, unique, range(0, SC_MAX_NAME_LENGTH)] wchar_t* lpDatabaseName,
//	[in] DWORD dwDesiredAccess,
//	[out] LPSC_RPC_HANDLE lpScHandle
//	);
//lpMachineName：一种 SVCCTL_HANDLEW（第 2.2.3 节）数据类型，它定义指向以空字符结尾的 UNICODE 字符串的指针，该字符串指定服务器的机器名称。
//lpDatabaseName：指向以空结尾的 UNICODE 字符串的指针，该字符串指定要打开的 SCM 数据库的名称。该参数必须设置为 NULL、“ServicesActive”或“ServicesFailed”。
//dwDesiredAccess：一个值，指定对数据库的访问。这必须是第 3.1.4 节中指定的值之一。
//客户端还必须具有 SC_MANAGER_CONNECT 访问权限。
//lpScHandle：一种 LPSC_RPC_HANDLE 数据类型，用于定义新打开的 SCM 数据库的句柄。
func (s *Session) NewSMB2OpenSCManagerWRequest(treeId uint32, fileId []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	machinename := string(util.Random(6)) + "\x00"
	databaseName := "ServicesActive\x00"
	buffer := OpenSCManagerWStruct{
		MachineName: machineName{
			MaxCount:    uint32(len(machinename)),
			ActualCount: uint32(len(machinename)),
			MachineName: encoder.ToUnicode(machinename),
		},
		Database: database{
			MaxCount:    uint32(len(databaseName)),
			ActualCount: uint32(len(databaseName)),
			Database:    encoder.ToUnicode(databaseName),
		},
		AccessMask: SC_MANAGER_CREATE_SERVICE | SC_MANAGER_CONNECT,
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             1,
		ContextId:          0,
		OpNum:              ROpenSCManagerW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewSMB2OpenSCManagerWResponse() OpenSCManagerWResponse {
	return OpenSCManagerWResponse{
		ContextHandle: make([]byte, 20),
	}
}

// 打开服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6d0a4225-451b-4132-894d-7cef7aecfd2d
type ROpenServiceWRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"` //OpenSCManagerW 句柄
	ServiceName   serviceName
	//Reserved      uint8
	AccessMask uint32
}

type ROpenServiceWResponseStruct struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

func (s *Session) NewSMB2ROpenServiceWRequest(treeId uint32, fileId, contextHandle []byte, servicename string) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	serName := servicename + "\x00"
	buffer := ROpenServiceWRequestStruct{
		ContextHandle: contextHandle,
		ServiceName: serviceName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			ServiceName: encoder.ToUnicode(serName),
		},
		AccessMask: SERVICE_ALL_ACCESS,
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             2,
		ContextId:          0,
		OpNum:              ROpenServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewSMB2ROpenServiceWResponse() ROpenServiceWResponseStruct {
	return ROpenServiceWResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}

// 创建服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6a8ca926-9477-4dd4-b766-692fab07227e
type RCreateServiceWRequestStruct struct {
	ContextHandle       []byte `smb:"fixed:20"` //OpenSCManagerW 句柄
	ServiceName         serviceName
	DisplayName         displayName
	AccessMask          uint32
	ServiceType         uint32
	ServiceStartType    uint32
	ServiceErrorControl uint32
	BinaryPathName      binaryPathName
	NULLPointer         uint32
	TagId               uint32
	NULLPointer2        uint32
	DependSize          uint32
	NULLPointer3        uint32
	NULLPointer4        uint32
	PasswordSize        uint32
}

type serviceName struct {
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32
	ServiceName []byte
	Reserved    uint16
}

type displayName struct {
	ReferentId  uint32 `smb:"offset:DisplayName"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32
	DisplayName []byte
	Reserved    uint16
}

type binaryPathName struct {
	MaxCount       uint32
	Offset         uint32
	ActualCount    uint32
	BinaryPathName []byte
}

// RCreateServiceW响应结构
type RCreateServiceWResponseStruct struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	TagId              uint32
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6a8ca926-9477-4dd4-b766-692fab07227e
// dwServiceType 类型
const (
	SERVICE_KERNEL_DRIVER       = 0x00000001
	SERVICE_FILE_SYSTEM_DRIVER  = 0x00000002
	SERVICE_WIN32_OWN_PROCESS   = 0x00000010
	SERVICE_WIN32_SHARE_PROCESS = 0x00000020
	SERVICE_INTERACTIVE_PROCESS = 0x00000100
)

// dwStartType类型
const (
	SERVICE_BOOT_START   = 0x00000000
	SERVICE_SYSTEM_START = 0x00000001
	SERVICE_AUTO_START   = 0x00000002
	SERVICE_DEMAND_START = 0x00000003
	SERVICE_DISABLED     = 0x00000004
)

// dwErrorControl类型
const (
	SERVICE_ERROR_IGNORE   = 0x00000000
	SERVICE_ERROR_NORMAL   = 0x00000001
	SERVICE_ERROR_SEVERE   = 0x00000002
	SERVICE_ERROR_CRITICAL = 0x00000003
)

func (s *Session) NewSMB2RCreateServiceWRequest(treeId uint32, fileId, contextHandle []byte, servicename, uploadPathFile string) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	serName := servicename + "\x00"
	uploadpathFile := uploadPathFile + "\x00"
	buffer := RCreateServiceWRequestStruct{
		ContextHandle: contextHandle,
		ServiceName: serviceName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			ServiceName: encoder.ToUnicode(serName),
		},
		DisplayName: displayName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			DisplayName: encoder.ToUnicode(serName),
		},
		AccessMask:          SERVICE_ALL_ACCESS,
		ServiceType:         SERVICE_WIN32_OWN_PROCESS,
		ServiceStartType:    SERVICE_DEMAND_START,
		ServiceErrorControl: SERVICE_ERROR_IGNORE,
		BinaryPathName: binaryPathName{
			MaxCount:       uint32(len(uploadpathFile)),
			ActualCount:    uint32(len(uploadpathFile)),
			BinaryPathName: encoder.ToUnicode(uploadpathFile),
		},
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             3,
		ContextId:          0,
		OpNum:              RCreateServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewSMB2RCreateServiceWResponse() RCreateServiceWResponseStruct {
	return RCreateServiceWResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}

// 启动服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/d9be95a2-cf01-4bdc-b30f-6fe4b37ada16
type RStartServiceWRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"` //20字节，创建服务返回的句柄
	Argc          uint32 //argv字符串数量
	Argv          []byte `smb:"fixed:4"` //4字节，unicode对象指针
}

type RStartServiceWResponseStruct struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	StubData           uint32
}

// 启动服务封装
func (s *Session) NewSMB2RStartServiceWRequest(treeId uint32, fileId, contextHandle []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	buffer := RStartServiceWRequestStruct{
		ContextHandle: contextHandle,
		Argc:          0,
		Argv:          encoder.ToUnicode("0"),
	}
	fragLength := 24 + util.SizeOfStruct(buffer)
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             4,
		ContextId:          0,
		OpNum:              RStartServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

// 启动服务响应封装
func NewSMB2RStartServiceWResponse() RStartServiceWResponseStruct {
	return RStartServiceWResponseStruct{}
}

// 关闭服务句柄
type RCloseServiceHandleRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"`
}

type RCloseServiceHandleResponseStruct struct {
	SMB2ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

func (s *Session) NewSMB2RCloseServiceHandleRequest(treeId uint32, fileId, contextHandle []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = s.messageId
	pduHeader.SMB2Header.SessionId = s.sessionId
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	buffer := RCloseServiceHandleRequestStruct{ContextHandle: contextHandle}
	fragLength := 24 + util.SizeOfStruct(buffer)
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             6,
		ContextId:          0,
		OpNum:              RStartServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewSMB2RCloseServiceHandleResponse() RCloseServiceHandleResponseStruct {
	return RCloseServiceHandleResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}

// 服务安装
func (s *Session) ServiceInstall(servicename string, uploadPathFile string) (service string, err error) {
	var fileId []byte
	// 建立ipc$管道
	treeId, err := s.SMB2TreeConnect("IPC$")
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	r := SMB2CreateRequestStruct{
		OpLock:             SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: Impersonation,
		AccessMask:         FILE_OPEN_IF,
		FileAttributes:     FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        FILE_SHARE_READ,
		CreateDisposition:  FILE_OPEN_IF,
		CreateOptions:      FILE_NON_DIRECTORY_FILE,
	}
	fileId, err = s.SMB2CreateRequest(treeId, "svcctl", r)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	// 绑定svcctl函数
	err = s.SMB2PDUBind(treeId, fileId, ms.NTSVCS_UUID, ms.NTSVCS_VERSION)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	req := s.NewSMB2OpenSCManagerWRequest(treeId, fileId)
	_, err = s.send(req)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	s.Debug("Read svcctl response", nil)
	req1 := s.NewSMB2ReadRequest(treeId, fileId)
	buf, err1 := s.send(req1)
	if err1 != nil {
		s.Debug("", err1)
		return "", err
	}
	res := NewSMB2OpenSCManagerWResponse()
	s.Debug("Unmarshalling OpenSCManagerW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to OpenSCManagerW service active to " + ms.StatusMap[res.SMB2Header.Status])
	}
	s.Debug("Completed OpenSCManagerW ", nil)
	// 获取OpenSCManagerW句柄
	contextHandle := res.ContextHandle
	// 打开服务
	s.Debug("Sending svcctl OpenServiceW request", nil)
	req2 := s.NewSMB2ROpenServiceWRequest(treeId, fileId, contextHandle, servicename)
	buf, err = s.send(req2)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	s.Debug("Read svcctl OpenServiceW response", nil)
	reqRead := s.NewSMB2ReadRequest(treeId, fileId)
	buf, err = s.send(reqRead)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	res1 := NewSMB2ROpenServiceWResponse()
	s.Debug("Unmarshalling ROpenServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res1); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	//if res.SMB2Header.Status != ms.STATUS_SUCCESS {
	//	return errors.New("Failed to ROpenServiceW to " + ms.StatusMap[res.SMB2Header.Status])
	//}
	s.Debug("Completed ROpenServiceW ", nil)
	// 创建服务
	s.Debug("Sending svcctl RCreateServiceW request", nil)
	// uploadPathFile %systemroot%\xxx
	req3 := s.NewSMB2RCreateServiceWRequest(treeId, fileId, contextHandle, servicename, uploadPathFile)
	buf, err = s.send(req3)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	s.Debug("Read svcctl RCreateServiceW response", nil)
	reqRead = s.NewSMB2ReadRequest(treeId, fileId)
	buf, err = s.send(reqRead)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	res2 := NewSMB2RCreateServiceWResponse()
	s.Debug("Unmarshalling RCreateServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res2); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	s.Debug("Completed RCreateServiceW to ["+servicename+"] ", nil)
	// 得到创建服务后的服务句柄
	serviceHandle := res2.ContextHandle
	// 启动服务
	// bug: 服务启动失败,0x20，原因未知
	s.Debug("Sending svcctl RStartServiceW request", nil)
	req4 := s.NewSMB2RStartServiceWRequest(treeId, fileId, serviceHandle)
	buf, err = s.send(req4)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	s.Debug("Read svcctl RStartServiceW response", nil)
	reqRead = s.NewSMB2ReadRequest(treeId, fileId)
	buf, err = s.send(reqRead)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	res3 := NewSMB2RStartServiceWResponse()
	s.Debug("Unmarshalling RStartServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res3); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res3.SMB2Header.Status != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to RStartServiceW to " + ms.StatusMap[res3.SMB2Header.Status])
	}
	s.Debug("Completed RStartServiceW ", nil)
	// 关闭服务管理句柄
	s.Debug("Sending svcctl RCloseServiceHandle request", nil)
	req5 := s.NewSMB2RCloseServiceHandleRequest(treeId, fileId, serviceHandle)
	buf, err = s.send(req5)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	s.Debug("Read svcctl RCloseServiceHandle response", nil)
	reqRead = s.NewSMB2ReadRequest(treeId, fileId)
	buf, err = s.send(reqRead)
	if err != nil {
		s.Debug("", err)
		return "", err
	}
	res4 := NewSMB2RCloseServiceHandleResponse()
	s.Debug("Unmarshalling RCloseServiceHandle response", nil)
	if err = encoder.Unmarshal(buf, &res3); err != nil {
		s.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res4.ReturnCode != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to RCloseServiceHandle to " + ms.StatusMap[res4.ReturnCode])
	}
	s.Debug("Completed RCloseServiceHandle ", nil)
	return servicename, nil
}
