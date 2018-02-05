package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"os"
)

const (
	name        = "grender"
	version     = "0.0.1"
	description = "grender for 3d models rendering."
)

var (
	config *Config = nil
)

func start_gin() error {
	gin.SetMode(*config.Server.Mode)

	var r *gin.Engine = gin.New()

	if *config.Server.Mode == gin.ReleaseMode {
		r.Use(SeeLogger())
	} else {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.POST("/render3d", render3d)

	return r.Run(fmt.Sprintf("%s:%d", *config.Server.ListenIP, *config.Server.ListenPort)) // listen and serve on 0.0.0.0:1323
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = description

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start the grender",
			Flags: []cli.Flag{
				/*
					cli.BoolFlag{
						Name:  "daemon, d",
						Usage: "make a daemon",
					},
				*/
				cli.StringFlag{
					Name:  "config-file, c",
					Usage: "attach config file",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Usage: "set gin mode", // debug, test or release
				},
			},
			Action: func(ctx *cli.Context) error {
				conf, err := NewConfig(ctx.String("config-file"))
				if err != nil {
					return cli.NewExitError(err.Error(), 86)
				}
				config = conf

				mode := ctx.String("mode")
				debug := true
				if *config.Server.Mode != gin.ReleaseMode {
					debug = false
				}
				if mode == gin.DebugMode {
					*config.Server.Mode = mode
					debug = true
				} else if mode == gin.TestMode {
					*config.Server.Mode = mode
					debug = false
				} else if mode == gin.ReleaseMode {
					*config.Server.Mode = mode
					debug = false
				}

				logFile := *config.Log.LogFile
				accessFile := *config.Log.AccessFile
				maxSize := *config.Log.MaxSize
				maxRolls := *config.Log.MaxRolls
				if err = InitLogger(logFile, maxSize, maxRolls, debug); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				if err = InitAccessLogger(accessFile, maxSize, maxRolls); err != nil {
					return cli.NewExitError(err.Error(), 2)
				}

				return start_gin()
			},
		},
		/*
			{
				Name:  "stop",
				Usage: "fast shutdown",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "reload",
				Usage: "reloading the configuration file",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
		*/
	}

	app.Run(os.Args)
}
