package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Server struct {
		Host       string `toml:"host"`
		Port       int    `toml:"port"`
		StaticPath string `toml:"static_path"`
	} `toml:"server"`

	Compiler struct {
		Compiler       string   `toml:"compiler"`
		CompileFlags   []string `toml:"compile_flags"`
		CompileTimeout int      `toml:"compile_timeout"`
		RunTimeout     int      `toml:"run_timeout"`
		TempDir        string   `toml:"temp_dir"`
	} `toml:"compiler"`

	Auth struct {
		UsersFile      string `toml:"users_file"`
		SessionTimeout int    `toml:"session_timeout"`
	} `toml:"auth"`

	WebSocket struct {
		Path         string `toml:"path"`
		PingInterval int    `toml:"ping_interval"`
		BufferSize   int    `toml:"buffer_size"`
	} `toml:"websocket"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	if _, err := toml.DecodeFile(configPath, &AppConfig); err != nil {
		return err
	}
	log.Printf("配置加载成功: %s", configPath)
	return nil
}
