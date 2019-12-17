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
	templates = &Template{
		templates: template.Must(template.New(``).Delims(`[[`, `]]`).Funcs(
			template.FuncMap{"dir": dir},
		).
			ParseGlob(
				filepath.Join(cfg.Templates, "*.html"))),
	}
}
