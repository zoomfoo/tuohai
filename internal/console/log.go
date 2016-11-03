package console

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
)

type Console struct {
	access_log, error_log log.Logger
}

var StdLog *Console

func New(l log.Logger) *Console {
	access_log, error_log := l, l
	// Prefix
	access_log.SetOutput(CreateLogFile(fmt.Sprintf("../../log/%saccess.log", l.Prefix())))
	error_log.SetOutput(CreateLogFile(fmt.Sprintf("../../log/%serror.log", l.Prefix())))
	StdLog = &Console{access_log: access_log, error_log: error_log}
	return StdLog
}

func CreateLogFile(path string) *os.File {
	r, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return r
}

func (c *Console) Error(err error) {
	c.error_log.Output(2, fmt.Sprintf("ERROR: %v \n", err))
}

func (c *Console) Errorf(format string, v ...interface{}) {
	c.error_log.Output(2, fmt.Sprintf(format, v))
}

func (c *Console) Warning() {

}

func (c *Console) Info() {

}

func (c *Console) Log(arg ...interface{}) {
	c.access_log.Output(2, fmt.Sprintf("Log: %v\n", arg))
}

func (c *Console) Logf(format string, v ...interface{}) {
	c.access_log.Output(2, fmt.Sprintf(format, v))
}

func (c *Console) Debug() {

}

func Debugf() {

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		StdLog.Log(method, statusCode, clientIP, path, latency, comment)
	}
}
