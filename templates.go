// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"html/template"
	"path"
	"path/filepath"
)

type Template struct {
	templates *template.Template
}

var (
	templates *Template
)

func dir(s string) string {
	return path.Dir(s)
}

func LoadTemplates() {
	tmpl := template.New(``)
	delims := len(cfg.Delim)
	if delims != 0 {
		if delims > 2 {
			tmpl = tmpl.Delims(cfg.Delim[:2], cfg.Delim[2:])
		}
	}

	templates = &Template{
		templates: template.Must(tmpl.Funcs(
			template.FuncMap{"dir": dir},
		).
			ParseGlob(
				filepath.Join(cfg.Templates, "*.html"))),
	}
}
