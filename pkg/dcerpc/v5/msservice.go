package v5

import (
	"encoding/hex"
	"errors"
	"go-impacket/pkg/encoder"
	"go-impacket/pkg/ms"
	"go-impacket/pkg/smb/smb2"
)

// 此文件提供访问windows服务管理安装/删除

// 服务安装
func (c *Client) ServiceInstall(servicename string, uploadPathFile string) (service string, err error) {
	var fileId []byte
	// 建立ipc$管道
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	r := smb2.CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_OPEN_IF,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_READ,
		CreateDisposition:  smb2.FILE_OPEN,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	fileId, err = c.CreateRequest(treeId, "svcctl", r)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	// 绑定svcctl函数
	err = c.PDUBind(treeId, fileId, ms.NTSVCS_UUID, ms.NTSVCS_VERSION)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	req := c.NewOpenSCManagerWRequest(treeId, fileId)
	_, err = c.Send(req)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	c.Debug("Read svcctl response", nil)
	req1 := c.NewReadRequest(treeId, fileId)
	buf, err1 := c.Send(req1)
	if err1 != nil {
		c.Debug("", err1)
		return "", err
	}
	res := NewOpenSCManagerWResponse()
	c.Debug("Unmarshalling OpenSCManagerW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to OpenSCManagerW service active to " + ms.StatusMap[res.SMB2Header.Status])
	}
	c.Debug("Completed OpenSCManagerW ", nil)
	// 获取OpenSCManagerW句柄
	contextHandle := res.ContextHandle
	// 打开服务
	c.Debug("Sending svcctl OpenServiceW request", nil)
	req2 := c.NewROpenServiceWRequest(treeId, fileId, contextHandle, servicename)
	buf, err = c.Send(req2)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	c.Debug("Read svcctl OpenServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	res1 := NewROpenServiceWResponse()
	c.Debug("Unmarshalling ROpenServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res1); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	//if res.SMB2Header.Status != ms.STATUS_SUCCESS {
	//	return errors.New("Failed to ROpenServiceW to " + ms.StatusMap[res.SMB2Header.Status])
	//}
	c.Debug("Completed ROpenServiceW ", nil)
	// 创建服务
	c.Debug("Sending svcctl RCreateServiceW request", nil)
	// uploadPathFile %systemroot%\xxx
	req3 := c.NewRCreateServiceWRequest(treeId, fileId, contextHandle, servicename, uploadPathFile)
	buf, err = c.Send(req3)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	c.Debug("Read svcctl RCreateServiceW response", nil)
	reqRead = c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	res2 := NewRCreateServiceWResponse()
	c.Debug("Unmarshalling RCreateServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res2); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	c.Debug("Completed RCreateServiceW to ["+servicename+"] ", nil)
	// 得到创建服务后的服务句柄
	serviceHandle := res2.ContextHandle
	// 启动服务
	c.Debug("Sending svcctl RStartServiceW request", nil)
	req4 := c.NewRStartServiceWRequest(treeId, fileId, serviceHandle)
	buf, err = c.Send(req4)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	c.Debug("Read svcctl RStartServiceW response", nil)
	reqRead = c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	res3 := NewRStartServiceWResponse()
	c.Debug("Unmarshalling RStartServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res3); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res3.SMB2Header.Status != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to RStartServiceW to " + ms.StatusMap[res3.SMB2Header.Status])
	}
	c.Debug("Completed RStartServiceW ", nil)
	// 关闭服务管理句柄
	c.Debug("Sending svcctl RCloseServiceHandle request", nil)
	req5 := c.NewRCloseServiceHandleRequest(treeId, fileId, serviceHandle)
	buf, err = c.Send(req5)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	c.Debug("Read svcctl RCloseServiceHandle response", nil)
	reqRead = c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return "", err
	}
	res4 := NewRCloseServiceHandleResponse()
	c.Debug("Unmarshalling RCloseServiceHandle response", nil)
	if err = encoder.Unmarshal(buf, &res3); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res4.ReturnCode != ms.STATUS_SUCCESS {
		return "", errors.New("Failed to RCloseServiceHandle to " + ms.StatusMap[res4.ReturnCode])
	}
	c.Debug("Completed RCloseServiceHandle ", nil)
	return servicename, nil
}
