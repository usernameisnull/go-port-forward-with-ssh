package config

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
)

type Config struct {
	Pem            string `ini:"pem"`
	PemPass        string `ini:"pem_pass"`
	JumpServer     string `ini:"jump_server"`
	Local          string `ini:"local"`
	Remote         string `ini:"remote"`
	JumpServerUser string `ini:"jump_server_user"`
	LogPath        string `ini:"log_path"`
}

func ConfigParse() (c Config) {
	path := getCurrentRunningDir()
	f, err := ini.Load(filepath.Join(path, "config.ini"))
	if err != nil {
		panic(err)
	}
	s := f.Section("Default")
	if err := s.StrictMapTo(&c); err != nil {
		panic(err)
	}
	return
}

func getCurrentRunningDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
