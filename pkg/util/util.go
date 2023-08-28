package util

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// 提供一些常用方法

// 计算结构体大小
func SizeOfStruct(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	var sum int
	switch v.Kind() {
	case reflect.Map:
		sum = 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Slice, reflect.Array:
		sum = 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.String:
		sum = 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Struct:
		sum = 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		sum = int(v.Type().Size())
	case reflect.Interface:
		if !v.IsNil() {
			return sizeof(reflect.ValueOf(v.Interface()))
		}
	default:
		return 0
	}
	return sum
}

// 读文件
func ReadFile(filename string) ([]byte, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// 处理PDU uuid 转成字节数组
func PDUUuidFromBytes(uuid string) []byte {
	s := strings.ReplaceAll(uuid, "-", "")
	b, _ := hex.DecodeString(s)
	r := []byte{b[3], b[2], b[1], b[0], b[5], b[4], b[7], b[6], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}
	return r
}

func Random(n int) []byte {
	const alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alpha[b%byte(len(alpha))]
	}
	return bytes
}

func DealCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// 计算IP地址范围的起始和结束
	startIP := ip.Mask(ipnet.Mask)
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)
	for i := range endIP {
		endIP[i] |= ^ipnet.Mask[i]
	}

	var ips []string
	// 逐个增加IP地址
	for ip := startIP; bytes.Compare(ip, endIP) <= 0; ip = nextIP(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

// 下一个IP地址
func nextIP(ip net.IP) net.IP {
	nextIP := make(net.IP, len(ip))
	copy(nextIP, ip)

	for i := len(nextIP) - 1; i >= 0; i-- {
		nextIP[i]++
		if nextIP[i] > 0 {
			break
		}
	}

	return nextIP
}

func DealAsterisk(s string) ([]string, error) {
	i := strings.Count(s, "*")
	switch i {
	case 1:
		return DealCIDR(strings.Replace(s, "*", "0", -1) + "/24")
	case 2:
		return DealCIDR(strings.Replace(s, "*", "0", -1) + "/16")
	case 3:
		return DealCIDR(strings.Replace(s, "*", "0", -1) + "/8")
	default:
		return nil, errors.New("invalid Asterisk format")
	}
}

func DealHyphen(s string) ([]string, error) {
	tmp := strings.Split(s, ".")
	if len(tmp) != 4 {
		return nil, errors.New("invalid IP range format")
	}
	iprange_tmp := strings.Split(tmp[3], "-")
	if len(iprange_tmp) != 2 {
		return nil, errors.New("invalid IP range format")
	}
	head, err := strconv.Atoi(iprange_tmp[0])
	if err != nil {
		return nil, err
	}
	tail, err := strconv.Atoi(iprange_tmp[1])
	if err != nil {
		return nil, err
	}
	if head < 0 || head > 255 || tail < 0 || tail > 255 || head > tail {
		return nil, errors.New("invalid IP range")
	}
	var ips []string
	for i := head; i <= tail; i++ {
		ips = append(ips, tmp[0]+"."+tmp[1]+"."+tmp[2]+"."+strconv.Itoa(i))
	}
	return ips, nil
}

func IpParse(s string) ([]string, error) {
	ipStrings := strings.Split(strings.Trim(s, ","), ",")
	var ips []string
	for i := 0; i < len(ipStrings); i++ {
		if strings.Contains(ipStrings[i], "*") {
			// 192.168.0.*
			ips_tmp, err := DealAsterisk(ipStrings[i])
			if err != nil {
				return nil, err
			}
			ips = append(ips, ips_tmp...)
		} else if strings.Contains(ipStrings[i], "/") {
			// 192.168.0.1/24
			ips_tmp, err := DealCIDR(ipStrings[i])
			if err != nil {
				return nil, err
			}
			ips = append(ips, ips_tmp...)
		} else if strings.Contains(ipStrings[i], "-") {
			// 192.668.0.1-255
			ips_tmp, err := DealHyphen(ipStrings[i])
			if err != nil {
				return nil, err
			}
			ips = append(ips, ips_tmp...)
		} else {
			// single ip
			if net.ParseIP(ipStrings[i]) != nil {
				ips = append(ips, ipStrings[i])
			} else {
				return nil, fmt.Errorf("invalid IP address: %s", ipStrings[i])
			}
		}
	}
	return ips, nil
}
