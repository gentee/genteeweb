// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/parser"
)

type Render struct {
	Content template.HTML
	Title   string
	Logo    *Logo
	Menu    []*MenuItem
	Nav     []*NavItem
}

var (
	ErrNotFound = errors.New(`Not found`)
	ErrContent  = errors.New(`Invalid content`)
)

func RenderPage(url string) (string, error) {
	var (
		err    error
		ok     bool
		page   *Page
		render Render
	)
	if cfg.mode == ModeLive {
		mutexPage.RLock()
	}
	if page, ok = pages[url]; !ok {
		if !strings.HasSuffix(url, `.html`) {
			if page, ok = pages[path.Join(url, `index.html`)]; !ok {
				page = pages[path.Join(url, `readme.html`)]
			}
		}
	}
	if cfg.mode == ModeLive {
		mutexPage.RUnlock()
	}
	if page == nil {
		return ``, ErrNotFound
	}
	file := filepath.Join(cfg.WebDir, filepath.FromSlash(page.url))
	var exist bool
	if cfg.mode != ModeDynamic {
		if _, err := os.Stat(file); err == nil {
			exist = true
		}
	}
	switch cfg.mode {
	case ModeLive:
	case ModeCache:
	case ModeStatic:
		if !exist {
			return ``, ErrNotFound
		}
	}
	if exist {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return ``, err
		}
		return string(data), nil
	}
	if len(page.Template) == 0 {
		page.Template = page.parent.Template
	}
	tpl := page.Template
	if len(tpl) == 0 {
		return page.body, err
	}
	buf := bytes.NewBuffer([]byte{})

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	var markDown bytes.Buffer
	if err = markdown.Convert([]byte(page.body), &markDown); err != nil {
		return ``, err
	}
	render.Content = template.HTML(markDown.String())
	render.Title = page.Title
	render.Logo = page.parent.Logo
	render.Menu = page.parent.Menu
	render.Nav = page.parent.Nav
	if err = templates.templates.ExecuteTemplate(buf, tpl+`.html`, render); err != nil {
		return ``, err
	}
	if cfg.mode != ModeDynamic {
		err = ioutil.WriteFile(file, buf.Bytes(), os.ModePerm)
	}
	return buf.String(), err
}
