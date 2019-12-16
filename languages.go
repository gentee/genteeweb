// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"path"
	"strings"
)

var (
	langNative = map[string]string{
		`de`: `Deutsch`,
		`en`: `English`,
		`es`: `Español`,
		`fr`: `Français`,
		`ja`: `日本語 (にほんご)`,
		`it`: `Italiano`,
		`ko`: `한국어`,
		`pt`: `Português`,
		`ru`: `Русский`,
	}
)

func getPattern(parent *Content, pattern *[]string) (ret string) {
	if parent.parent != nil {
		ret = getPattern(parent.parent, pattern)
	}
	if len(parent.lang) != 0 {
		ret = parent.lang
		*pattern = append(*pattern, `%s`)
	} else {
		*pattern = append(*pattern, parent.dir)
	}
	return ret
}

func LangList(page *Page) []*MenuItem {
	var (
		ret     []*MenuItem
		pattern []string
		curCode string
	)
	if len(page.parent.Langs) < 2 {
		return nil
	}
	curCode = getPattern(page.parent, &pattern)
	patternUrl := `/` + path.Join(pattern...) + `/` + path.Base(page.url)
	if len(curCode) == 0 {
		return nil
	}
	ret = append(ret, &MenuItem{
		Title: langNative[curCode],
	})
	for _, code := range page.parent.Langs {
		if curCode == code {
			continue
		}
		url := fmt.Sprintf(patternUrl, code)
		if code == page.parent.DefDir {
			url = strings.ReplaceAll(url, `/`+code, ``)
		}
		if pageLang, ok := pages[url]; ok {
			ret = append(ret, &MenuItem{
				Title: langNative[code],
				Href:  pageLang.url,
			})
		}
	}
	if len(ret) > 1 {
		return ret
	}
	return nil
}
