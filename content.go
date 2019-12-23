// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
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
	Langs    []string          `yaml:"langs"`
	DefDir   string            `yaml:"defdir"`
	GitHub   string            `yaml:"github"`

	dir      string
	url      string
	lang     string // language code
	parent   *Content
	children []*Content
	pages    []*Page
}

var (
	pages      = map[string]*Page{}
	redirs     = map[string]string{}
	contents   = &Content{}
	reTitle, _ = regexp.Compile(`\s*#\s*(.*)`)
	redir      string
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

func FileToURL(page *Page) string {
	path := path.Join(page.parent.url, filepath.Base(page.file))
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
	if len(page.Title) == 0 {
		list := reTitle.FindStringSubmatch(body)
		if len(list) > 1 {
			page.Title = list[1]
		}
	}
	page.url = FileToURL(page)
	page.body = body
	pages[page.url] = page
	if len(redir) > 0 {
		redirUrl := strings.Replace(page.url, path.Dir(redir), redir, 1)
		redirs[redirUrl] = page.url
	}
	return page
}

func readDir(dpath string, owner *Content) {
	files, err := ioutil.ReadDir(dpath)
	if err != nil {
		golog.Fatal(err)
	}
	var (
		dirs []string
	)
	for _, file := range files {
		fname := file.Name()
		if file.IsDir() {
			dirs = append(dirs, filepath.Join(dpath, fname))
		} else {
			page := ReadContent(filepath.Join(dpath, fname), owner)
			if page != nil {
				owner.pages = append(owner.pages, page)
			}
		}
	}
	readme := path.Join(owner.url, `readme.html`)
	index := path.Join(owner.url, `index.html`)
	if pages[readme] != nil && pages[index] == nil {
		page := pages[readme]
		page.url = index
		delete(pages, readme)
		pages[index] = page
		redirs[readme] = index
	}
	for _, item := range owner.Paths {
		dir, err := filepath.Abs(item.Dir)
		if err != nil {
			golog.Fatal(err)
		}
		dirs = append(dirs, dir)
		cfg.paths = append(cfg.paths, dir)
	}
	for _, dir := range dirs {
		var (
			url, lang string
		)
		base := strings.ToLower(filepath.Base(dir))
		url = path.Join(owner.url, base)
		if base == owner.DefDir {
			redir = url
			url = owner.url
		}
		for _, ilang := range owner.Langs {
			if ilang == base {
				lang = ilang
			}
		}
		child := &Content{
			dir:    base,
			parent: owner,
			url:    url,
			lang:   lang,
		}
		readDir(dir, child)
		redir = ``
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
				if len(lang) == 0 {
					child.Nav = parent.Nav
				} else {
					child.Nav = copyNav(parent.Nav)
					navMenu(child)
				}
			}
			if child.Langs == nil {
				child.Langs = parent.Langs
			}
			if len(child.DefDir) == 0 {
				child.DefDir = parent.DefDir
			}
			if len(child.GitHub) == 0 {
				child.GitHub = parent.GitHub
			}
			if len(child.Template) == 0 {
				child.Template = parent.Template
			}
			parent = parent.parent
		}
		owner.children = append(owner.children, child)
	}
	navMenu(owner)
}

func LoadContent() {
	golog.Info(`Reading content...`)
	os.Chdir(filepath.Dir(cfg.Content))
	contents.url = `/`
	readDir(cfg.Content, contents)
}

func UpdateContent(page *Page) error {
	curUrl := page.url
	page = ReadContent(page.file, page.parent)
	if page == nil {
		return ErrContent
	}
	page.url = curUrl
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

func walkNav(owner *Content, list []*NavItem) {
	for i, item := range list {
		if len(item.Children) > 0 {
			walkNav(owner, item.Children)
		}
		if len(item.Href) == 0 || strings.HasPrefix(item.Href, `/`) ||
			strings.IndexByte(item.Href, ':') > 0 {
			continue
		}
		if page, ok := pages[path.Join(owner.url, item.Href)+`.html`]; ok {
			list[i].Href = page.url
			if len(item.Title) == 0 {
				if len(page.TitleMenu) > 0 {
					list[i].Title = template.JS(page.TitleMenu)
				} else {
					list[i].Title = template.JS(page.Title)
				}
			}
		}
	}
}

func navMenu(owner *Content) {
	if len(owner.DefDir) == 0 {
		walkNav(owner, owner.Nav)
	}
}

func copyNav(nav []*NavItem) (ret []*NavItem) {
	for _, item := range nav {
		retNav := NavItem{
			Title: item.Title,
			Href:  item.Href,
		}
		if len(item.Children) > 0 {
			retNav.Children = copyNav(item.Children)
		}
		ret = append(ret, &retNav)
	}
	return
}
