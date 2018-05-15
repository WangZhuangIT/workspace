package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.xesv5.com/senior/lib/go/xeslog"
	"github.com/fsnotify/fsnotify"
)

type fileLoader struct {
	dir     string
	watcher *fsnotify.Watcher
}

func NewFileLoader(dir string) *fileLoader {
	l := new(fileLoader)
	l.dir = strings.TrimRight(dir, "/")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		xeslog.Fatal(err)
	}
	l.watcher = watcher
	err = l.watchDir()
	if err != nil {
		xeslog.Fatal(err)
	}
	return l
}

func (l *fileLoader) watchDir() error {
	err := l.watcher.Add(l.dir)
	if err != nil {
		return err
	}
	err = filepath.Walk(l.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err = l.watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func readDir(path string) ([]string, error) {
	filelist := []string{}
	var err error
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return filelist, err
	}
	for _, v := range files {
		if v.IsDir() {
			filelist_tmp, err := readDir(path + "/" + v.Name())
			if err != nil {
				return filelist, err
			}
			filelist = mergeSlice(filelist, filelist_tmp)
			continue
		}
		filelist = append(filelist, path+"/"+v.Name())
	}
	return filelist, err
}

func mergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func (l *fileLoader) Read() (map[string][]byte, error) {
	//遍历文件夹，设置文件的key，和val
	var err error
	data := make(map[string][]byte)
	files, err := readDir(l.dir)
	if err != nil {
		return data, err
	}
	for _, v := range files {
		data[strings.Replace(v, l.dir+"/", "", -1)], err = ioutil.ReadFile(v)
		if err != nil {
			return data, err
		}

	}
	return data, err
}

func (l *fileLoader) Watch(onChange func(map[string][]byte)) {
	for {
		select {
		case ev := <-l.watcher.Events:
			xeslog.Debug(ev)
			switch ev.Op {
			case fsnotify.Create:
				//判断创建的文件是否为文件夹，文件夹再次监听
				s, err := os.Stat(ev.Name)
				if err != nil {
					xeslog.Error(err)
				}
				if s.IsDir() {
					l.watcher.Add(ev.Name)
				}
			case fsnotify.Rename, fsnotify.Remove:
				l.watcher.Remove(ev.Name)
			}
			data, err := l.Read()
			if err != nil {
				xeslog.Error(err)
			} else {
				onChange(data)
			}

		case err := <-l.watcher.Errors:
			xeslog.Error(err)
		}
	}
}
