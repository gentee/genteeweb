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
)

type Render struct {
	Content template.HTML
}

var (
	ErrNotFound = errors.New(`Not found`)
	ErrContent  = errors.New(`Invalid content`)
)

func getTemplate(page *Page) string {
	if len(page.Template) > 0 {
		return page.Template
	}
	parent := page.parent
	for parent != nil {
		if len(parent.Template) > 0 {
			return parent.Template
		}
		parent = parent.parent
	}
	return ``
}

func RenderPage(url string) (string, error) {
	var (
		err    error
		ok     bool
		page   *Page
		render Render
	)
	if page, ok = pages[url]; !ok {
		if !strings.HasSuffix(url, `.html`) {
			if page, ok = pages[path.Join(url, `index.html`)]; !ok {
				page = pages[path.Join(url, `readme.html`)]
			}
		}
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
		if exist {
			if finfo, err := os.Stat(page.file); err == nil {
				if finfo.ModTime().After(page.modtime) {
					exist = false
					page = ReadContent(path.Dir(page.url), page.file, page.parent)
					if page == nil {
						return ``, ErrContent
					}
					pages[page.url] = page
					for i, cur := range page.parent.pages {
						if cur.url == page.url {
							page.parent.pages[i] = page
							break
						}
					}
				}
			}
		}
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
	tpl := getTemplate(page)
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
	)
	var markDown bytes.Buffer
	if err = markdown.Convert([]byte(page.body), &markDown); err != nil {
		return ``, err
	}
	render.Content = template.HTML(markDown.String())
	err = templates.templates.ExecuteTemplate(buf, tpl+`.html`, render)
	if cfg.mode != ModeDynamic {
		err = ioutil.WriteFile(file, buf.Bytes(), os.ModePerm)
	}
	return buf.String(), err
}
