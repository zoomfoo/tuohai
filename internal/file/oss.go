package file

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

func UploadFile(reader io.Reader) error {
	client, err := oss.New("http://img-cn-qingdao.aliyuncs.com", "muNWzl5jWgiNzDcq", "ixlGqqPQQxZzG8hZYIpqKs51o89qmB")
	if err != nil {
		return err
	}

	bucket, err := client.Bucket("zhizhiboom")
	if err != nil {
		return err
	}

	fd, err := os.Open("/Users/apple/Desktop/222.jpg")
	if err != nil {
		return err
	}
	defer fd.Close()

	err = bucket.PutObject("my-object.jpg", fd)
	if err != nil {
		return err
	}
	return nil
}
