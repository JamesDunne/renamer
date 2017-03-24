package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	//"strconv"
	//"strings"
	"regexp"
)

type FileInfoPath struct {
	FileInfo os.FileInfo
	Dir string
} 

func main() {
	localPath := `F:\Video\TV\The X-Files\`

	files := make([]FileInfoPath, 0, 50)
	dirs := make([]string, 0, 50)
	dirs = append(dirs, localPath)

	for (len(dirs) > 0) {
		// Open the directory to read its contents:
		dirPath := dirs[0]
		dirs = dirs[1:]
		df, dferr := os.Open(dirPath)
		if dferr != nil {
			log.Fatalln(dferr)
			return
		}

		// Read the directory entries:
		fis, fierr := df.Readdir(0)
		if fierr != nil {
			log.Fatalln(fierr)
			df.Close()
			return
		}
		df.Close()

		for _, fi := range fis {
			if fi.IsDir() {
				// Skip folders:
				if fi.Name() == "Extras" {
					continue
				}
				dirs = append(dirs, filepath.Join(dirPath, fi.Name()))
				continue
			}

			if filepath.Ext(fi.Name()) != ".avi" {
				continue
			}

			files = append(files, FileInfoPath{FileInfo: fi, Dir: dirPath})
		}
	}

	re := regexp.MustCompile(`^(.*) Season (\d\d) Episode (\d\d) - (.*)$`)
	for _, fpath := range files {
		oldName := fpath.FileInfo.Name()
		newName := re.ReplaceAllString(oldName, "${1} S${2}E${3} - ${4}")

		if (oldName == newName) {
			continue
		}

		// Rename file:
		fmt.Printf("%s\n", filepath.Join(fpath.Dir, newName))
		os.Rename(filepath.Join(fpath.Dir, oldName), filepath.Join(fpath.Dir, newName))
	}
}