package file

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	// "io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"tuohai/internal/uuid"
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

	filemd5 := FileName(*buf)
	fmt.Println("filemd5: ", filemd5)
	name := fmt.Sprintf("avatar/%s/%s/%s.%s", filemd5[0:2], filemd5[2:4], uuid.NewV4().String(), suffix)
	fmt.Println("name: ", name)
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
