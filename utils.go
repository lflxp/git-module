// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func ExecCommandString(dir, name, cmd string) (string, error) {
	pipeline := exec.Command("/bin/sh", "-c", fmt.Sprintf("%s %s", name, cmd))
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Dir = dir
	pipeline.Stdout = &out
	pipeline.Stderr = &stderr
	err := pipeline.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

// objectCache provides thread-safe cache operations.
// TODO(@unknwon): Use sync.Map once requires Go 1.13.
type objectCache struct {
	lock  sync.RWMutex
	cache map[string]interface{}
}

func newObjectCache() *objectCache {
	return &objectCache{
		cache: make(map[string]interface{}),
	}
}

func (oc *objectCache) Set(id string, obj interface{}) {
	oc.lock.Lock()
	defer oc.lock.Unlock()

	oc.cache[id] = obj
}

func (oc *objectCache) Get(id string) (interface{}, bool) {
	oc.lock.RLock()
	defer oc.lock.RUnlock()

	obj, has := oc.cache[id]
	return obj, has
}

// isDir returns true if given path is a directory,
// or returns false when it's a file or does not exist.
func isDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// isFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func isFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// isExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func concatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%v - %s", err, stderr)
}
