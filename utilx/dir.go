// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/9

package utilx

import (
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() (dir string, err error) {
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	return
}
