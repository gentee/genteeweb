// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

type Page struct {
	Title     string `yaml:"title"`
	Desc      string `yaml:"description"`
	TitleMenu string `yaml:"titlemenu"`
	Template  string `yaml:"template"`

	body   string
	url    string
	parent *Content
	//	url  string
}

type Content struct {
	Template string `yaml:"template"`

	dir      string
	parent   *Content
	children []*Content
}

var (
	pages    = map[string]*Page{}
	contents *Content
)

func readFile(www, name string, owner *Content) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		golog.Fatal(err)
	}
	fname := filepath.Base(name)
	if fname == ContentFile {
		if err = yaml.Unmarshal(data, owner); err != nil {
			golog.Fatal(err)
		}
		return
	}
	page := &Page{
		parent: owner,
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
	fmt.Println(name, page)
}

func readDir(path, www string, owner *Content) {
	if owner == nil {
		owner = &Content{}
	}
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
			readFile(www, filepath.Join(path, fname), owner)
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
