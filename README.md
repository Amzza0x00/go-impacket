Go-Impacket
===
-------
基于golang实现的impacket  
> 目前仅实现smb2、dce/rpc协议
-------
示例
-------
```shell
psexec 172.20.10.5 Administrator DESKTOP-3397AU79 32ed87bdb5fdc5e9cba88547376818d4 test.exe go-impacket/test/
```
参考
-------
基于以下代码修改，并修复pass the hash不成功问题  
**smb**: https://github.com/stacktitan/smb