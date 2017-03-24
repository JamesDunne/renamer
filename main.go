package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	//"strconv"
	//"strings"
	"regexp"
)

func main() {
	path := flag.String("path", ".", "local path")
	isRecursive := flag.Bool("recurse", false, "recurse subdirectories")
	matchRe := flag.String("match", "(.*)", "regexp match")
	replaceRe := flag.String("replace", "", "regexp replace")
	doReplace := flag.Bool("exec", false, "actually rename files")
	flag.Parse()

	// `^(.*) Season (\d\d) Episode (\d\d) - (.*)$`
	re := regexp.MustCompile(*matchRe)

	dirs := make([]string, 0, 50)
	dirs = append(dirs, *path)

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
				if !*isRecursive {
					continue
				}

				dirs = append(dirs, filepath.Join(dirPath, fi.Name()))
				continue
			}

			oldName := fi.Name()
			newName := oldName
			if *replaceRe == "" {
				if re.MatchString(oldName) {
					fmt.Printf("%s\n", filepath.Join(dirPath, oldName))
				}
				continue
			}

			// "${1} S${2}E${3} - ${4}"
			newName = re.ReplaceAllString(oldName, *replaceRe)
			if (oldName == newName) {
				continue
			}

			fmt.Printf("%s%c{%s -> %s}\n", dirPath, filepath.Separator, oldName, newName)
			if *doReplace {
				// Rename file:
				os.Rename(
					filepath.Join(dirPath, oldName),
					filepath.Join(dirPath, newName),
				)
			}
		}
	}
}
