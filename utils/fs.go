package utils

import (
	"fmt"
	"io/fs"
	"log"
)

// might be useful to know the folder structure of afero fs
func PrintAllFiles(fsys fs.FS, dir string) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		log.Printf("Error reading directory %s: %v\n", dir, err)
		return
	}

	for _, entry := range entries {
		path := dir + "/" + entry.Name()
		if dir == "." {
			path = entry.Name() // Avoid leading "./"
		}

		if entry.IsDir() {
			PrintAllFiles(fsys, path)
		} else {
			fmt.Println(path)
		}
	}
}
