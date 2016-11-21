package file

import (
	"io"

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

func UploadFile(reader io.Reader) *NetPath {
	client, err := oss.New(osshost, "muNWzl5jWgiNzDcq", "ixlGqqPQQxZzG8hZYIpqKs51o89qmB")
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	bucket, err := client.Bucket("zhizhiboom")
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	err = bucket.PutObject("my-object.jpg", reader)
	if err != nil {
		return &NetPath{P: "", E: err}
	}
	return &NetPath{P: "http://zhizhiboom.img-cn-qingdao.aliyuncs.com/my-object.jpg", E: nil}
}
