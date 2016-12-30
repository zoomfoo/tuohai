package pushnotify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type PNotify struct {
	tlsConfig *tls.Config
	submsg    chan redis.Message //订阅回来的管道(暂时。。。需要修改)
	exitChan  chan int
}

func NewPNotify() (*PNotify, error) {
	return &PNotify{
		submsg: make(chan redis.Message, 20),
	}, nil
}

//notify主入口
func (n *PNotify) Main(opts *Options) error {
	config, err := buildTLSConfig(opts)
	if err != nil || config == nil {
		return err
	}
	n.tlsConfig = config
	topic := "com.yunzujia.woke"

	//初始化redis连接
	InitRedis(opts.RedisHost, "")

	//从redis中订阅
	go Subscribers(n.submsg, opts.Subscribers)

	for {
		select {
		case sub := <-n.submsg:
			var api Notification
			err := json.Unmarshal(sub.Data, &api)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("从消息中间件中序列化数据错误: ", err)
			fmt.Println("反序列化后内容: ", api)
			fmt.Printf("为序列化json串: %s", sub.Data)

			for _, device := range api.Ios {
				if len(device.Token) > 1 {
					for _, tok := range device.Token {
						n.HTTP2Client(opts, &Handle{
							token: tok,
							alert: device.Alert,
							badge: device.Badge,
							sound: device.Sound,
							topic: topic,
						})
					}
				} else {
					n.HTTP2Client(opts, &Handle{
						token: device.Token[0],
						alert: device.Alert,
						badge: device.Badge,
						sound: device.Sound,
						topic: topic,
					})
				}
			}

		}
	}

	return nil
}

func buildTLSConfig(opts *Options) (tlsConfig *tls.Config, err error) {
	if opts.TLSCert != "" && opts.TLSKey != "" {
		return
	}

	if opts.P12File == "" {
		fmt.Println("P12File is empty")
		return
	}

	cert, err := FromP12File(opts.P12File, opts.P12Password)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   opts.TLSMinVersion,
		MaxVersion:   tls.VersionTLS12,
	}

	return
}
