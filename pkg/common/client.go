package common

// 客户端连接封装

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/Amzza0x00/go-impacket/pkg/encoder"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

// 会话结构
type Client struct {
	IsSigningRequired bool
	IsAuthenticated   bool
	debug             bool
	securityMode      uint16
	messageId         uint64
	sessionId         uint64
	conn              net.Conn
	dialect           uint16
	options           *ClientOptions
	trees             map[string]uint32
}

// 连接参数
type ClientOptions struct {
	Host        string
	Port        int
	Workstation string
	Domain      string
	User        string
	Password    string
	Hash        string
}

func (c *Client) Debug(msg string, err error) {
	if c.debug {
		log.Println("[ DEBUG ] ", msg)
		if err != nil {
			debug.PrintStack()
		}
	}
}

func (c *Client) SMBSend(req interface{}) (res []byte, err error) {
	buf, err := encoder.Marshal(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	b := new(bytes.Buffer)
	if err = binary.Write(b, binary.BigEndian, uint32(len(buf))); err != nil {
		c.Debug("", err)
		return
	}
	c.Debug("Raw:\n"+hex.Dump(append(b.Bytes(), buf...)), nil)
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	if _, err = rw.Write(append(b.Bytes(), buf...)); err != nil {
		c.Debug("", err)
		return
	}
	rw.Flush()
	var size uint32
	if err = binary.Read(rw, binary.BigEndian, &size); err != nil {
		c.Debug("", err)
		return
	}
	if size > 0x00FFFFFF {
		return nil, errors.New("Invalid NetBIOS Session message")
	}
	data := make([]byte, size)
	l, err := io.ReadFull(rw, data)
	if err != nil {
		c.Debug("", err)
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
	c.messageId++
	return data, nil
}

//func (c *Client) TCPSend(req interface{}) (res []byte, err error) {
//	buf, err := encoder.Marshal(req)
//	if err != nil {
//		c.Debug("", err)
//		return nil, err
//	}
//	c.Debug("Raw:\n"+hex.Dump(buf), nil)
//	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
//	if _, err = rw.Write(buf); err != nil {
//		c.Debug("", err)
//		return
//	}
//	rw.Flush()
//	data := make([]byte, 4096)
//	c.conn.Read(data)
//	c.messageId++
//	return data, nil
//}

func (c *Client) TCPSend(req interface{}) (res []byte, err error) {
	buf, err := encoder.Marshal(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	c.Debug("Raw:\n"+hex.Dump(buf), nil)
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	if _, err = rw.Write(buf); err != nil {
		c.Debug("", err)
		return nil, err
	}
	if err = rw.Flush(); err != nil {
		c.Debug("", err)
		return nil, err
	}

	var responseData bytes.Buffer
	responseBuffer := make([]byte, 4096)

	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 5)) // 设置读取操作的超时时间为5秒

		n, err := c.conn.Read(responseBuffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 超时，认为没有更多响应
				break
			} else if err == io.EOF {
				// 服务器关闭连接，认为没有更多响应
				break
			}
			c.Debug("", err)
			return nil, err
		}

		responseData.Write(responseBuffer[:n])

		// 根据结束条件判断是否继续接收响应
		// if shouldStop(responseData.Bytes()) {
		//     break
		// }

		// 动态调整缓冲区大小，扩展为当前大小的两倍
		if responseData.Len() >= len(responseBuffer) {
			newBuffer := make([]byte, responseData.Len()*2)
			copy(newBuffer, responseData.Bytes())
			responseBuffer = newBuffer
		}
	}

	c.messageId++
	return responseData.Bytes(), nil
}

func (c *Client) WithDebug(debug bool) *Client {
	c.debug = debug
	return c
}

func (c *Client) WithSecurityMode(securityMode uint16) *Client {
	c.securityMode = securityMode
	return c
}

func (c *Client) GetSecurityMode() uint16 {
	return c.securityMode
}

func (c *Client) GetMessageId() uint64 {
	return c.messageId
}

func (c *Client) WithSessionId(sessionId uint64) *Client {
	c.sessionId = sessionId
	return c
}

func (c *Client) GetSessionId() uint64 {
	return c.sessionId
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) WithConn(conn net.Conn) *Client {
	c.conn = conn
	return c
}

func (c *Client) WithDialect(dialect uint16) *Client {
	c.dialect = dialect
	return c
}

func (c *Client) WithOptions(clientOptions *ClientOptions) *Client {
	c.options = clientOptions
	return c
}

func (c *Client) GetOptions() *ClientOptions {
	return c.options
}

func (c *Client) WithTrees(trees map[string]uint32) *Client {
	c.trees = trees
	return c
}

func (c *Client) GetTrees() map[string]uint32 {
	return c.trees
}

func (c *Client) Close() error {
	if c.conn != nil {
		// 关闭连接之前，设置一个较短的读写超时时间，确保及时返回
		c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		c.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))

		// 关闭连接
		if err := c.conn.Close(); err != nil {
			c.Debug("", err)
			return err
		}
		c.conn = nil
	}
	return nil
}
