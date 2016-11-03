package http

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"tuohai/internal/console"
	"tuohai/internal/pb/IM_Message"
	"tuohai/internal/pb/MsgSrv"
)

func SendLogicMsg(ConnLogicRPCAddress string, p *IM_Message.IMMsgData) (*IM_Message.IMMsgDataAck, error) {
	conn, err := grpc.Dial(ConnLogicRPCAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Millisecond))
	if err != nil {
		console.StdLog.Error(err)
		log.Println(err)
		return nil, err
	}
	defer conn.Close()
	c := MsgSrv.NewMsgLogicClient(conn)

	r, err := c.SendMsg(context.Background(), p)
	if err != nil {
		log.Println(err)
		console.StdLog.Error(err)
		return r, err
	}
	return r, nil
}
