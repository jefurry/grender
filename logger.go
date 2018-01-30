package main

import (
	"bytes"
	"github.com/cihub/seelog"
	"text/template"
)

type SeeLogConfig struct {
	LogFile  string
	MaxSize  int
	MaxRolls int
	Debug    bool
}

var Logger seelog.LoggerInterface = nil

func GetLogger() seelog.LoggerInterface {
	return Logger
}

func InitLogger(logFile string, maxSize, maxRolls int, debug bool) error {
	if Logger != nil {
		Logger.Flush()
	}

	config := `
	<seelog>
		<outputs formatid="main">
			<filter levels="debug,info,warn,critical,error">
				{{if eq .Debug true}}
					<console />
				{{else}}
					<rollingfile type="size" filename="{{.LogFile}}" maxsize="{{.MaxSize}}" maxrolls="{{.MaxRolls}}" />
				{{end}}
			</filter>
		</outputs>
		<formats>
			<format id="main" format="%Date/%Time [%LEV] %Msg%n"/>
		</formats>
	</seelog>
	`

	var lc *SeeLogConfig = &SeeLogConfig{
		LogFile:  logFile,
		Debug:    debug,
		MaxSize:  maxSize,
		MaxRolls: maxRolls,
	}

	t, err := template.New("seelog").Parse(config)
	if err != nil {
		return err
	}

	out := bytes.NewBuffer(nil)
	err = t.Execute(out, lc)
	if err != nil {
		return err
	}

	Logger, err = seelog.LoggerFromConfigAsBytes(out.Bytes())
	if err != nil {
		return err
	}

	return nil
}
