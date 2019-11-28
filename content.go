// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

type Page struct {
	Title     string `yaml:"title"`
	Desc      string `yaml:"description"`
	TitleMenu string `yaml:"titlemenu"`
	Template  string `yaml:"template"`

	body    string
	url     string
	file    string
	modtime time.Time
	parent  *Content
}

type Content struct {
	Template string `yaml:"template"`

	dir      string
	parent   *Content
	children []*Content
	pages    []*Page
}

var (
	pages    = map[string]*Page{}
	contents = &Content{}
)

func ReadContent(www, name string, owner *Content) *Page {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		golog.Fatal(err)
	}
	fname := filepath.Base(name)
	if fname == ContentFile {
		if err = yaml.Unmarshal(data, owner); err != nil {
			golog.Fatal(err)
		}
		return nil
	}
	page := &Page{
		parent: owner,
		file:   name,
	}
	if cfg.mode == ModeCache || cfg.mode == ModeLive {
		finfo, err := os.Stat(name)
		if err != nil {
			golog.Fatal(err)
		}
		page.modtime = finfo.ModTime()
	}
	body := strings.TrimSpace(string(data))
	lenMD := len(MDHead)
	if strings.HasPrefix(body, MDHead) {
		if off := strings.Index(body[lenMD:], "\n"+MDHead); off >= 0 {
			head := strings.Trim(body[:off+lenMD], "-\r\n")
			body = body[off+2*lenMD+1:]
			if err = yaml.Unmarshal([]byte(head), &page); err != nil {
				golog.Fatal(err)
			}
		}
	}
	fname = fname[:len(fname)-len(filepath.Ext(fname))]
	page.url = strings.ToLower(path.Join(www, fname+`.html`))
	page.body = body
	pages[page.url] = page
	return page
}

func readDir(path, www string, owner *Content) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		golog.Fatal(err)
	}
	var dirs []string
	for _, file := range files {
		fname := file.Name()
		if file.IsDir() {
			dirs = append(dirs, fname)
		} else {
			page := ReadContent(www, filepath.Join(path, fname), owner)
			if page != nil {
				owner.pages = append(owner.pages, page)
			}
		}
	}
	for _, dir := range dirs {
		child := &Content{
			dir:    strings.ToLower(dir),
			parent: owner,
		}
		readDir(filepath.Join(path, dir), www+`/`+dir, child)
		owner.children = append(owner.children, child)
	}
}

func LoadContent() {
	golog.Info(`Reading content...`)
	readDir(cfg.Content, `/`, contents)
}
