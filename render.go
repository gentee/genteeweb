// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"path"
	"strings"
)

type Render struct {
	Content string
}

var (
	ErrNotFound = errors.New(`Not found`)
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
	tpl := getTemplate(page)
	if len(tpl) == 0 {
		return page.body, err
	}
	buf := bytes.NewBuffer([]byte{})
	render.Content = page.body
	err = templates.templates.ExecuteTemplate(buf, tpl+`.html`, render)
	return buf.String(), err
}
