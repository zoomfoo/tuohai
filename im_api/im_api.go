package im_api

import (
	"fmt"
	"net"
	"os"

	"tuohai/internal/svc"
	"tuohai/internal/util"
	"tuohai/models"
)

type ImApi struct {
	Opts         *Options
	httpListener net.Listener
	waitGroup    util.WaitGroupWrapper
}

func New(opts *Options) *ImApi {
	return &ImApi{Opts: opts}
}

func (api *ImApi) Main() {
	httpListener, err := net.Listen("tcp", api.Opts.HTTPAddress)
	if err != nil {
		os.Exit(1)
	}

	models.InitDB(api.Opts.MysqlOptions())
	models.InitRedis(api.Opts.RedisHost, api.Opts.RedisPwd)

	api.httpListener = httpListener
	fmt.Println("LISTEN: ", api.httpListener.Addr().String())
	fmt.Println("PID: ", api.Opts.ID)
	api.waitGroup.Wrap(func() {
		svc.HttpService(api.httpListener, newHTTPServer())
	})
}
