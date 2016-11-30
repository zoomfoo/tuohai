package file_api

import (
	"fmt"
	"net"
	"os"

	"tuohai/file_api/models"
	"tuohai/internal/svc"
	"tuohai/internal/util"
)

type FileApi struct {
	Opts         *Options
	httpListener net.Listener
	waitGroup    util.WaitGroupWrapper
}

func New(opts *Options) *FileApi {
	return &FileApi{Opts: opts}
}

func (file *FileApi) Main() {
	httpListener, err := net.Listen("tcp", Opts.HTTPAddress)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	models.InitDB(file.Opts.MysqlOptions())
	models.InitRedis(file.Opts.RedisHost, file.Opts.RedisPwd)

	file.httpListener = httpListener
	fmt.Println("LISTEN: ", file.httpListener.Addr().String())
	fmt.Println("PID: ", file.Opts.ID)
	file.waitGroup.Wrap(func() {
		svc.HttpService(file.httpListener, newHTTPServer())
	})
}
