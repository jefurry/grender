package main

import (
	"bytes"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

var (
	accessLogger     seelog.LoggerInterface = nil
	ByteBufferWriter *bytes.Buffer          = bytes.NewBuffer(nil)
)

func SeeLogger() gin.HandlerFunc {
	f := gin.LoggerWithWriter(ByteBufferWriter)

	return func(c *gin.Context) {
		defer accessLogger.Flush()
		f(c)
		msg := ByteBufferWriter.String()
		accessLogger.Info(msg)
	}
}

func InitAccessLogger(accessFile string, maxSize, maxRolls int) error {
	if accessLogger != nil {
		return nil
	}

	config := fmt.Sprintf(`
	<seelog>
		<outputs formatid="main">
			<filter levels="info">
				<rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
			</filter>
		</outputs>
		<formats>
			<format id="main" format="%s"/>
		</formats>
	</seelog>
	`, accessFile, maxSize, maxRolls, "[%LEV] %Msg%n")

	l, err := seelog.LoggerFromConfigAsBytes([]byte(config))
	if err != nil {
		return err
	}
	accessLogger = l

	return nil
}
