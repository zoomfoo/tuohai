package open_api

import (
	"fmt"
	"net"
	"os"

	"tuohai/internal/svc"
	"tuohai/internal/util"
	"tuohai/models"
)

type OpenApi struct {
	Opts         *Options
	httpListener net.Listener
	waitGroup    util.WaitGroupWrapper
}

func NewOpenApi(opts *Options) *OpenApi {
	return &OpenApi{Opts: opts}
}

func (api *OpenApi) Main() {
	httpListener, err := net.Listen("tcp", api.Opts.HTTPAddress)
	if err != nil {
		os.Exit(1)
	}

	models.InitDB(api.Opts.MysqlOptions())

	api.httpListener = httpListener
	fmt.Println("LISTEN: ", api.httpListener.Addr().String())
	fmt.Println("PID: ", api.Opts.ID)
	api.waitGroup.Wrap(func() {
		svc.HttpService(api.httpListener, newHTTPServer())
	})
}

func (oa *OpenApi) Close() error {
	if oa == nil {
		return fmt.Errorf("%s", "OpenApi can't be empty")
	}
	oa.waitGroup.Wait()
	return nil
}
