package httpReq

import (
	"crypto/tls"
	"os"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

func NewClient() *req.Client {
	client := req.C().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	}).SetTimeout(1 * time.Minute)
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		client = client.DevMode()
	}
	return client
}

func NewFeiShuClient() *req.Client {
	client := req.C().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	}).SetTimeout(1 * time.Minute)
	return client
}
