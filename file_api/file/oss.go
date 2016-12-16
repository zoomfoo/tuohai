package file

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	// "io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"tuohai/file_api/options"
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
	filemd5 := FileName(buf.Bytes())
	fmt.Println("filemd5: ", filemd5)
	name := fmt.Sprintf("file/%s/%s/%s.%s", filemd5[0:2], filemd5[2:4], uuid.NewV4().String(), suffix)
	fmt.Println("name: ", name)

	np := upload(options.Opts.FileOSSHost,
		options.Opts.AccessKeyId,
		options.Opts.AccessKeySecret,
		options.Opts.FileBucket,
		suffix, name, buf)
	np.P = options.Opts.FileHost + np.P
	return np
}

func AvatarUpload(suffix string, buf *bytes.Buffer) *NetPath {
	filemd5 := FileName(buf.Bytes())
	fmt.Println("filemd5: ", filemd5)
	name := fmt.Sprintf("avatar/%s/%s/%s.%s", filemd5[0:2], filemd5[2:4], uuid.NewV4().String(), suffix)
	fmt.Println("name: ", name)

	return upload(options.Opts.OSSHost,
		options.Opts.AccessKeyId,
		options.Opts.AccessKeySecret,
		options.Opts.AvatarBucket,
		suffix, name, buf)
}

func upload(OSSHost, AccessKeyId, AccessKeySecret, Bucket, suffix, name string, buf *bytes.Buffer) *NetPath {
	client, err := oss.New(OSSHost, AccessKeyId, AccessKeySecret)
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	bucket, err := client.Bucket(Bucket)
	if err != nil {
		return &NetPath{P: "", E: err}
	}

	err = bucket.PutObject(name, buf)
	if err != nil {
		return &NetPath{P: "", E: err}
	}
	return &NetPath{P: "/" + name, E: nil}
}

func FileName(d []byte) string {
	h := md5.New()
	h.Write(d)
	return hex.EncodeToString(h.Sum(nil))
}
