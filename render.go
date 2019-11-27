// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	buf := bytes.NewBuffer([]byte{})
	err := t.templates.ExecuteTemplate(buf, name, data)

	_, err = buf.WriteTo(w)
	return err
}

func LoadTemplates() *Template {
	t := &Template{
		templates: template.Must(template.New(``).Delims(`[[`, `]]`).ParseGlob(
			filepath.Join(cfg.Root, "templates", "*.html"))),
	}
	return t
}
