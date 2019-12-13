// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"html/template"
	"path/filepath"
)

type Template struct {
	templates *template.Template
}

var (
	templates *Template
)

func toJS(s string) template.JS {
	return template.JS(s)
}

func LoadTemplates() {
	templates = &Template{
		templates: template.Must(template.New(``).Delims(`[[`, `]]`).Funcs(
			template.FuncMap{"toJS": toJS},
		).
			ParseGlob(
				filepath.Join(cfg.Templates, "*.html"))),
	}
}
