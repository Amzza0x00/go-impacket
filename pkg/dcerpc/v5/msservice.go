package v5

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go-impacket/pkg/dcerpc"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb/smb2"
	"go-impacket/pkg/util"
	"strings"
)

// 此文件提供访问windows服务管理安装/删除

// smb->上传文件，返回文件名
func (c *SMBClient) FileUpload(file, Path string) (filename string, err error) {
	treeId, err := c.TreeConnect("C$")
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	createRequestStruct := smb2.CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_CREATE,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_WRITE,
		CreateDisposition:  smb2.FILE_OVERWRITE_IF,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	// 文件名不能超过11位，如果超过则随机生成
	var newFilename string
	if len(file) <= 11 {
		newFilename = file
	} else {
		// 切分拿到扩展名
		fileInfo := strings.Split(file, ".")
		newFilename = string(util.Random(7)) + "." + fileInfo[len(fileInfo)-1]
	}
	fileId, err := c.CreateRequest(treeId, newFilename, createRequestStruct)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	err = c.WriteRequest(treeId, Path, file, fileId)
	if err != nil {
		c.Debug("", err)
		return newFilename, err
	}
	// 关闭目录连接
	c.TreeDisconnect("C$")
	return newFilename, nil
}

// smb->打开scm，返回scm服务句柄
func (c *SMBClient) OpenSvcManager(treeId, callId uint32) (fileid, handler []byte, err error) {
	createRequestStruct := smb2.CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_OPEN_IF,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_READ,
		CreateDisposition:  smb2.FILE_OPEN,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	fileId, err := c.CreateRequest(treeId, "svcctl", createRequestStruct)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	// 绑定svcctl函数
	err = c.MSRPCBind(treeId, fileId, ms.NTSVCS_UUID, ms.NTSVCS_VERSION)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	openSvcManagerRequest := NewOpenSCManagerWRequest()
	openSvcManagerRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, openSvcManagerRequest)
	_, err = c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	c.Debug("Read OpenSCManagerW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err := c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	smbRes := smb2.NewReadResponse()
	res := NewOpenSCManagerWResponse()
	c.Debug("Unmarshalling OpenSCManagerW response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.ReturnCode]) == 0 {
			msg := fmt.Sprintf("Failed to OpenSCManagerW service active code : 0x%08x", res.ReturnCode)
			return nil, nil, errors.New(msg)
		} else {
			return nil, nil, errors.New("Failed to OpenSCManagerW service active : " + dcerpc.RpcStatusCodes[res.ReturnCode])
		}
	}
	c.Debug("Completed OpenSCManagerW ", nil)
	// 获取OpenSCManagerW句柄
	contextHandle := res.ContextHandle
	return fileId, contextHandle, nil
}

// smb->打开服务
func (c *SMBClient) OpenService(treeId uint32, fileId, contextHandle []byte, servicename string, callId uint32) (err error) {
	// 打开服务
	c.Debug("Sending svcctl OpenServiceW request", nil)
	rOpenServiceRequest := NewROpenServiceWRequest(contextHandle, servicename)
	rOpenServiceRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, rOpenServiceRequest)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl OpenServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	smbRes := smb2.NewReadResponse()
	res := NewROpenServiceWResponse()
	c.Debug("Unmarshalling ROpenServiceW response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.ReturnCode]) == 0 {
			msg := fmt.Sprintf("Failed to ROpenServiceW service active code : 0x%08x", res.ReturnCode)
			return errors.New(msg)
		} else {
			return errors.New("Failed to ROpenServiceW service active : " + dcerpc.RpcStatusCodes[res.ReturnCode])
		}
	}
	c.Debug("Completed ROpenServiceW ", nil)
	return nil
}

// smb->创建服务，返回创建服务后的实例句柄
func (c *SMBClient) CreateService(treeId uint32, fileId, contextHandle []byte, servicename, uploadPathFile string, callId uint32) (handler []byte, err error) {
	// 创建服务
	c.Debug("Sending svcctl RCreateServiceW request", nil)
	rCreateServiceWRequest := NewRCreateServiceWRequest(contextHandle, servicename, uploadPathFile)
	rCreateServiceWRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, rCreateServiceWRequest)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	c.Debug("Read svcctl RCreateServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	smbRes := smb2.NewReadResponse()
	res := NewRCreateServiceWResponse()
	c.Debug("Unmarshalling RCreateServiceW response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.ReturnCode]) == 0 {
			msg := fmt.Sprintf("Failed to RCreateServiceW service active code : 0x%08x", res.ReturnCode)
			return nil, errors.New(msg)
		} else {
			return nil, errors.New("Failed to RCreateServiceW service active : " + dcerpc.RpcStatusCodes[res.ReturnCode])
		}
	}
	c.Debug("Completed RCreateServiceW to ["+servicename+"] ", nil)
	// 得到创建服务后的服务句柄
	serviceHandle := res.ContextHandle
	return serviceHandle, nil
}

// smb->启动服务
func (c *SMBClient) StartService(treeId uint32, fileId, serviceHandle []byte, callId uint32) (err error) {
	// 启动服务
	c.Debug("Sending svcctl RStartServiceW request", nil)
	rStartServiceWRequest := NewRStartServiceWRequest(serviceHandle)
	rStartServiceWRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, rStartServiceWRequest)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl RStartServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	smbRes := smb2.NewReadResponse()
	res := NewRStartServiceWResponse()
	c.Debug("Unmarshalling RStartServiceW response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.StubData != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.StubData]) == 0 {
			msg := fmt.Sprintf("Failed to RStartServiceW service active code : 0x%08x", res.StubData)
			return errors.New(msg)
		} else {
			return errors.New("Failed to RStartServiceW service active : " + dcerpc.RpcStatusCodes[res.StubData])
		}
	}
	c.Debug("Completed RStartServiceW ", nil)
	return nil
}

// smb->删除服务
func (c *SMBClient) DeleteService(treeId uint32, fileId, serviceHandle []byte, callId uint32) (err error) {
	c.Debug("Sending svcctl RDeleteService request", nil)
	rDeleteServiceRequest := NewRDeleteServiceRequest(serviceHandle)
	rDeleteServiceRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, rDeleteServiceRequest)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl RDeleteService response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	smbRes := smb2.NewReadResponse()
	res := NewRDeleteServiceResponse()
	c.Debug("Unmarshalling RDeleteService response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.ReturnCode]) == 0 {
			msg := fmt.Sprintf("Failed to RDeleteService service active code : 0x%08x", res.ReturnCode)
			return errors.New(msg)
		} else {
			return errors.New("Failed to RDeleteService service active : " + dcerpc.RpcStatusCodes[res.ReturnCode])
		}
	}
	c.Debug("Completed RDeleteService ", nil)
	return nil
}

// smb->关闭scm句柄
func (c *SMBClient) CloseService(treeId uint32, fileId, serviceHandle []byte, callId uint32) error {
	// 关闭服务管理句柄
	c.Debug("Sending svcctl RCloseServiceHandle request", nil)
	rCloseServiceHandleRequest := NewRCloseServiceHandleRequest(serviceHandle)
	rCloseServiceHandleRequest.CallId = callId
	req := c.NewWriteRequest(treeId, fileId, rCloseServiceHandleRequest)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl RCloseServiceHandle response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	smbRes := smb2.NewReadResponse()
	res := NewRCloseServiceHandleResponse()
	c.Debug("Unmarshalling RCloseServiceHandle response", nil)
	if err = encoder.Unmarshal(buf, &smbRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	// 切开smb头
	startIndex := len(buf) - int(smbRes.BlobLength)
	if err = encoder.Unmarshal(buf[startIndex:], &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != dcerpc.RPC_S_OK {
		if len(dcerpc.RpcStatusCodes[res.ReturnCode]) == 0 {
			msg := fmt.Sprintf("Failed to RCloseServiceHandle service active code : 0x%08x", res.ReturnCode)
			return errors.New(msg)
		} else {
			return errors.New("Failed to RCloseServiceHandle service active : " + dcerpc.RpcStatusCodes[res.ReturnCode])
		}
	}
	c.Debug("Completed RCloseServiceHandle ", nil)
	return nil
}

// 服务安装
func (c *SMBClient) ServiceInstall(servicename, file, path string) (service string, servicehandle []byte, err error) {
	// 上传文件
	filename, err := c.FileUpload(file, path)
	if err != nil {
		fmt.Println("[-]", err)
		return "", nil, err
	}
	//建立ipc$管道
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		fmt.Println("[-]", err)
		return "", nil, err
	}
	var callId uint32
	callId = 2
	// 打开服务管理
	svcctlFileId, svcctlHandler, err := c.OpenSvcManager(treeId, callId)
	callId++
	if err != nil {
		fmt.Println("[-]", err)
		return "", nil, err
	}
	// 打开服务
	err = c.OpenService(treeId, svcctlFileId, svcctlHandler, servicename, callId)
	if err != nil {
		fmt.Println("[-]", err)
		//return "", err
	}
	callId++
	// 创建服务
	uploadFilePath := "%systemdrive%\\" + filename
	serviceHandle, err := c.CreateService(treeId, svcctlFileId, svcctlHandler, servicename, uploadFilePath, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return "", nil, err
	}
	callId++
	// 启动服务
	err = c.StartService(treeId, svcctlFileId, serviceHandle, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return servicename, serviceHandle, err
	}
	callId++
	// 关闭服务管理
	err = c.CloseService(treeId, svcctlFileId, svcctlHandler, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return servicename, serviceHandle, err
	}
	return servicename, serviceHandle, nil
}

// 服务删除
func (c *SMBClient) ServiceDelete(serviceHandle []byte) (err error) {
	//建立ipc$管道
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		fmt.Println("[-]", err)
		return err
	}
	var callId uint32
	callId = 7
	// 打开服务管理
	svcctlFileId, svcctlHandler, err := c.OpenSvcManager(treeId, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return err
	}
	callId++
	// 删除服务
	err = c.DeleteService(treeId, svcctlFileId, serviceHandle, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return err
	}
	callId++
	// 关闭服务管理
	err = c.CloseService(treeId, svcctlFileId, svcctlHandler, callId)
	if err != nil {
		fmt.Println("[-]", err)
		return err
	}
	fmt.Println("[+] Service has been removed")
	return nil
}
