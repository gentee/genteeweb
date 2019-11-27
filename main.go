// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"github.com/kataras/golog"
)

func main() {
	golog.SetTimeFormat("2006/01/02 15:04:05")
	LoadSettings()
	LoadContent()
	RunServer()
}
