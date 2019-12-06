// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
)

var (
	watcher   *fsnotify.Watcher
	mutexPage = &sync.RWMutex{}
)

func StartLive() {
	var (
		err    error
		isPage bool
		page   *Page
		timer  *time.Timer
		event  fsnotify.Event
		ok     bool
	)
	eventFunc := func() {
		if strings.HasPrefix(event.Name, cfg.Templates) &&
			event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
			LoadTemplates()
			RemoveCache()
			fmt.Println("templ:", event, err)
			if err != nil {
				golog.Error(err)
			}
			return
		}
		fmt.Println("event:", event, len(watcher.Events))
		if event.Op&fsnotify.Write == fsnotify.Write {
			if strings.HasPrefix(event.Name, cfg.Content) {
				if filepath.Base(event.Name) == ContentFile {

				} else {
					isPage = true
				}
			}
		}
		if event.Op&fsnotify.Create == fsnotify.Create {
			if strings.HasPrefix(event.Name, cfg.Content) {
				/*					if filepath.Base(event.Name) == ContentFile {

									} else {
										isPage = true
									}*/
			}
		}
		if isPage {
			mutexPage.Lock()
			page = pages[FileToURL(event.Name)]
			if page != nil {
				if err := UpdateContent(page); err != nil {
					golog.Error(err)
				}
			}
			mutexPage.Unlock()
			isPage = false
		}
	}
	for {
		select {
		case event, ok = <-watcher.Events:
			if !ok {
				return
			}
			if filepath.Base(event.Name)[0] == '.' {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(time.Second, eventFunc)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			golog.Errorf("watcher error: %v", err)
		}
	}
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}

func WatchDirs() error {
	for _, dir := range []string{cfg.Content, cfg.Templates} {
		if err := filepath.Walk(dir, watchDir); err != nil {
			return err
		}
	}
	return nil
}
