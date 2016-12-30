package pushnotify

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

//使用http2 客户端想APNs 发送消息
func (n *PNotify) HTTP2Client(opts *Options, handle APNsapi) {
	var host string
	if opts.Production {
		host = ProductionServer
	} else {
		host = DevelopmentServer
	}

	client := NewClient(host, &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: n.tlsConfig,
		},
	})

	fmt.Println(string(handle.Aps()))

	header := &Header{ApnsTopic: handle.Topic()}

	res, err := client.Do(handle.Token(), header, handle.Aps())
	if err != nil {
		fmt.Println(res)
		switch err.(type) {
		case *ErrorResponse:
			fmt.Println("error response:", err.(*ErrorResponse).Reason)
		default:
			fmt.Println(err)
			return
		}
	}
	fmt.Println(res)
}

func FromP12File(filename string, password string) (tls.Certificate, error) {
	p12bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return FromP12Bytes(p12bytes, password)
}

func FromP12Bytes(bytes []byte, password string) (tls.Certificate, error) {
	key, cert, err := pkcs12.Decode(bytes, password)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}, nil
}
