// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

type Page struct {
	Title     string `yaml:"title"`
	Desc      string `yaml:"desc"`
	TitleMenu string `yaml:"titlemenu"`
	Template  string `yaml:"template"`
	Body      string
}

type Content struct {
	Template string `yaml:"template"`

	parent   *Content
	children []*Content
}

func LoadContent() {

}
