// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"html/template"
	"io/ioutil"
	"os"
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

type Logo struct {
	Name string `yaml:"name"`
}

type PathItem struct {
	Dir string `yaml:"dir"`
}

type MenuItem struct {
	Title string `yaml:"title"`
	Href  string `yaml:"href"`
}

type NavItem struct {
	Title    template.JS `yaml:"title"`
	Href     string      `yaml:"href"`
	Children []*NavItem  `yaml:"children"`
}

type Content struct {
	Template string            `yaml:"template"`
	Logo     *Logo             `yaml:"logo"`
	Params   map[string]string `yaml:"params"`
	Paths    []*PathItem       `yaml:"paths"`
	Menu     []*MenuItem       `yaml:"menu"`
	Nav      []*NavItem        `yaml:"nav"`

	dir      string
	parent   *Content
	children []*Content
	pages    []*Page
}

var (
	pages    = map[string]*Page{}
	contents = &Content{}
)

func RemoveCache() {
	files, err := ioutil.ReadDir(cfg.WebDir)
	if err != nil {
		golog.Error(err)
	}
	for _, file := range files {
		fname := file.Name()
		if fname == `assets` {
			continue
		}
		fname = filepath.Join(cfg.WebDir, fname)
		if file.IsDir() {
			err = os.RemoveAll(fname)
		} else {
			err = os.Remove(fname)
		}
		if err != nil {
			golog.Error(err)
		}
	}
	createDir(contents, cfg.WebDir, false)
}

func FileToURL(fName string) string {
	path := filepath.ToSlash(fName[len(cfg.Content):])
	if path[0] != '/' {
		path = `/` + path
	}
	path = path[:len(path)-len(filepath.Ext(path))]
	return strings.ToLower(path + `.html`)
}

func ReadContent(name string, owner *Content) *Page {
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
	page.url = FileToURL(name)
	page.body = body
	pages[page.url] = page
	return page
}

func readDir(path string, owner *Content) {
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
			page := ReadContent(filepath.Join(path, fname), owner)
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
		readDir(filepath.Join(path, dir), child)
		parent := child.parent
		for parent != nil {
			if child.Logo == nil {
				child.Logo = parent.Logo
			}
			if child.Params == nil {
				child.Params = parent.Params
			}
			if child.Menu == nil {
				child.Menu = parent.Menu
			}
			if child.Nav == nil {
				child.Nav = parent.Nav
			}
			if len(child.Template) == 0 {
				child.Template = parent.Template
			}
			parent = parent.parent
		}
		owner.children = append(owner.children, child)
	}
}

func LoadContent() {
	golog.Info(`Reading content...`)
	readDir(cfg.Content, contents)
}

func UpdateContent(page *Page) error {
	page = ReadContent(page.file, page.parent)
	if page == nil {
		return ErrContent
	}
	pages[page.url] = page
	for i, cur := range page.parent.pages {
		if cur.url == page.url {
			page.parent.pages[i] = page
			break
		}
	}
	os.Remove(filepath.Join(cfg.WebDir, filepath.FromSlash(page.url)))
	return nil
}
