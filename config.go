// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

const (
	ModeLive   = iota // compare the time of .md and .html
	ModeCache         // returns .html if it exists
	ModeStatic        // creates all .hmtl at the start
)

type Config struct {
	Root    string `yaml:"root"`    // Root directory. If it is empty - dir of yaml file
	Domain  string `yaml:"domain"`  // domain
	Port    int    `yaml:"port"`    // if empty, then 8080
	Content string `yaml:"content"` // content directory. content if it is empty
	LogDir  string `yaml:"logdir"`  // directory for log files, log if it is empty
	WebDir  string `yaml:"webdir"`  // directory for static web files, www if it is empty
	Mode    string `yaml:"mode"`    // mode: static, cache, live. Default, live
	Lang    string `yaml:"lang"`    // default language. By default, en
	mode    int
}

var (
	cfg Config
)

func LoadSettings() {
	var (
		ok      bool
		err     error
		cfgData []byte
	)
	cfgFile := DefaultCfgName
	if len(os.Args) > 1 {
		cfgFile = os.Args[1]
	} else if _, err := os.Stat("/etc/" + DefaultCfgName); err == nil {
		cfgFile = "/etc/" + DefaultCfgName
	}
	if cfgFile, err = filepath.Abs(cfgFile); err != nil {
		golog.Fatal(err)
	}
	if cfgData, err = ioutil.ReadFile(cfgFile); err != nil {
		golog.Fatal(err)
	}
	if err = yaml.Unmarshal(cfgData, &cfg); err != nil {
		golog.Fatal(err)
	}
	if len(cfg.Root) == 0 {
		cfg.Root = filepath.Dir(cfgFile)
	} else if cfg.Root, err = filepath.Abs(cfg.Root); err != nil {
		golog.Fatal(err)
	}
	if cfg.Port == 0 {
		cfg.Port = DefaultPort
	}
	if len(cfg.Mode) == 0 {
		cfg.Mode = DefaultMode
	}
	if len(cfg.Content) == 0 {
		cfg.Content = filepath.Join(cfg.Root, `content`)
	}
	if len(cfg.LogDir) == 0 {
		cfg.LogDir = filepath.Join(cfg.Root, `log`)
	}
	if len(cfg.WebDir) == 0 {
		cfg.WebDir = filepath.Join(cfg.Root, `www`)
	}
	if len(cfg.Lang) == 0 {
		cfg.Lang = DefaultLang
	}
	if cfg.mode, ok = map[string]int{
		`live`:   ModeLive,
		`cache`:  ModeCache,
		`static`: ModeStatic,
	}[cfg.Mode]; !ok {
		golog.Fatalf(`Unknown mode %s`, cfg.Mode)
	}
	fmt.Println(`CFG`, cfg)
}
