package v5

import (
	"go-impacket/pkg/smb/smb2"
)

type Client struct {
	smb2.Client
}

func SMBTransport() (client *Client, err error) {
	return &Client{}, nil
}
