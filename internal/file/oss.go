package file

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	// "fmt"
	// "io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var osshost = "http://img-cn-qingdao.aliyuncs.com"

type NetPath struct {
	P string //path
	E error  //error
}

func (p *NetPath) Path() (string, error) {
	return p.P, p.E
}

func UploadFile(suffix string, buf *bytes.Buffer) *NetPath {
	client, err := oss.New(osshost, "muNWzl5jWgiNzDcq", "ixlGqqPQQxZzG8hZYIpqKs51o89qmB")
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	bucket, err := client.Bucket("zhizhiboom")
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	name := FileName(*buf) + suffix

	err = bucket.PutObject(name, buf)
	if err != nil {
		return &NetPath{P: "", E: err}
	}
	return &NetPath{P: "http://zhizhiboom.img-cn-qingdao.aliyuncs.com/" + name, E: nil}
}

func FileName(buf bytes.Buffer) string {
	h := md5.New()
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}
