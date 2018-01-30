package main

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
)

var (
	// server
	DefaultListenIP   string = "0.0.0.0"
	DefaultListenPort int    = 1323
	DefaultMode       string = "release"

	// log
	DefaultLogFile    string = "grender.log"
	DefaultAccessFile string = "grender_access.log"
	DefaultMaxSize    int    = 10000000
	DefaultMaxRolls   int    = 30

	// render
	DefaultStlFilePaths   []string = []string{"examples/"}
	DefaultImageFilePaths []string = []string{"examples/images/"}

	// jwt
	DefaultExpires  int64  = 5
	DefaultIssuer   string = "3d@grender"
	DefaultAudience string = "urn:grender"
	DefaultSubject  string = "grender"
)

// jwt config
type JwtConfig struct {
	Expires  *int64  `json:"expires"`
	Issuer   *string `json:"issuer"`
	Audience *string `json:"audience"`
	Subject  *string `json:"subject"`
}

func NewJwtConfig() *JwtConfig {
	return &JwtConfig{}
}

func (jc *JwtConfig) SetDefault() error {
	if jc.Expires == nil {
		jc.Expires = &DefaultExpires
	}
	if jc.Issuer == nil {
		jc.Issuer = &DefaultIssuer
	}
	if jc.Audience == nil {
		jc.Audience = &DefaultAudience
	}
	if jc.Subject == nil {
		jc.Subject = &DefaultSubject
	}

	return nil
}

// server config
type ServerConfig struct {
	ListenIP   *string `json:"listen-ip"`
	ListenPort *int    `json:"listen-port"`
	Mode       *string `json:"mode"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (sc *ServerConfig) SetDefault() error {
	if sc.ListenIP == nil {
		sc.ListenIP = &DefaultListenIP
	}
	if sc.ListenPort == nil {
		sc.ListenPort = &DefaultListenPort
	}

	if sc.Mode == nil {
		sc.Mode = &DefaultMode
	}

	return nil
}

// log config
type LogConfig struct {
	LogFile    *string `json:"log-file"`
	AccessFile *string `json:"access-file"`
	MaxSize    *int    `json:"max-size"`
	MaxRolls   *int    `json:"max-rolls"`
}

func NewLogConfig() *LogConfig {
	return &LogConfig{}
}

func (lc *LogConfig) SetDefault() error {
	if lc.LogFile == nil {
		lc.LogFile = &DefaultLogFile
	}
	if err := EnsurePath(*lc.LogFile, false); err != nil {
		return err
	}

	if lc.AccessFile == nil {
		lc.AccessFile = &DefaultAccessFile
	}

	if lc.MaxSize == nil {
		lc.MaxSize = &DefaultMaxSize
	}
	if lc.MaxRolls == nil {
		lc.MaxRolls = &DefaultMaxRolls
	}

	return nil
}

// render config
type RenderConfig struct {
	StlFilePaths   []string `json:"stl-file-paths"`
	ImageFilePaths []string `json:"image-file-paths"`
}

func NewRenderConfig() *RenderConfig {
	return &RenderConfig{}
}

func (rc RenderConfig) SetDefault() error {
	if rc.StlFilePaths == nil {
		rc.StlFilePaths = DefaultStlFilePaths
	}
	if rc.ImageFilePaths == nil {
		rc.ImageFilePaths = DefaultImageFilePaths
	}

	for i, v := range rc.StlFilePaths {
		rc.StlFilePaths[i] = strings.TrimRight(v, "/") + "/"
	}
	for i, v := range rc.ImageFilePaths {
		rc.ImageFilePaths[i] = strings.TrimRight(v, "/") + "/"
	}

	return nil
}

// config
type Config struct {
	Server *ServerConfig `json:"server"`
	Log    *LogConfig    `json:"log"`
	Render *RenderConfig `json:"render"`
	Jwt    *JwtConfig    `json:"jwt"`
}

func NewConfig(configFile string) (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("config file: '%s' is not exists", configFile))
		} else {
			return nil, errors.New(fmt.Sprintf("config file: '%s'", err.Error()))
		}
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("config file: '%s'", err.Error()))
	}

	return cfg.SetDefault()
}

func (c *Config) SetDefault() (*Config, error) {
	var err error

	// server
	if c.Server == nil {
		c.Server = NewServerConfig()
	}
	if err = c.Server.SetDefault(); err != nil {
		return nil, err
	}

	// log
	if c.Log == nil {
		c.Log = NewLogConfig()
	}
	if err = c.Log.SetDefault(); err != nil {
		return nil, err
	}

	// render
	if c.Render == nil {
		c.Render = NewRenderConfig()
	}
	if err = c.Render.SetDefault(); err != nil {
		return nil, err
	}

	// jwt
	if c.Jwt == nil {
		c.Jwt = NewJwtConfig()
	}
	if err = c.Jwt.SetDefault(); err != nil {
		return nil, err
	}

	return c, nil
}
