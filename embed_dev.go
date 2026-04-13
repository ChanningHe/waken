//go:build dev

package main

import "io/fs"

func getFrontendFS() fs.FS {
	return nil
}
