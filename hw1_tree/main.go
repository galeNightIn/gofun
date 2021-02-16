package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)
type Map map[bool]string

type Bools []bool

type Files struct {
	files      []os.FileInfo
	filesCount int
}


func toStringBytes(size int64) string {
	if size == 0 {
		return "(empty)"
	}
	return strings.Join([]string{"(", strconv.FormatInt(size, 10), "b)"}, "")
}


func sortDirByName(slice []os.FileInfo, reversed bool) {
	lessFunc := func(i, j int) bool {
		result := slice[i].Name() < slice[j].Name()

		if reversed {
			result = !result
		}

		return result
	}

	sort.SliceStable(slice, lessFunc)
}


func getSortedFiles(path string, printFiles bool) (error, Files) {
	file, err := os.Open(path)

	if err != nil {
		return err, Files{}
	}

	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	fileInfoSlice, err := file.Readdir(0)
	if err != nil {
		return err, Files{}
	}

	// Oh shit holy fuck I want filter here
	if !printFiles {
		newSlice := make([]os.FileInfo, 0, 0)
		for _, file := range fileInfoSlice {
			if file.IsDir() {
				newSlice = append(newSlice, file)
			}
		}

		fileInfoSlice = newSlice
	}

	// We are good humans, let's sort our stuff
	sortDirByName(fileInfoSlice, false)
	filesCount := len(fileInfoSlice)

	return err, Files{fileInfoSlice, filesCount}

}

func subTree(out io.Writer, path string, printFiles bool, folderPattern Map, filePattern Map, parentDirs *Bools) (int, int, error) {
	err, files := getSortedFiles(path, printFiles)
	if err != nil { return 0, 0, err }

	dirCount, fileCount := 0, 0
	filesCount := files.filesCount

	for i, file := range files.files {
		count := i + 1
		name := file.Name()

		if !file.IsDir() {
			name += " " + toStringBytes(file.Size())
		}

		subPrint(out, count, filesCount, name, *parentDirs, folderPattern, filePattern)

		if file.IsDir() {
			parentDirs := append(*parentDirs, count == filesCount)

			newPath := strings.Join([]string{path, file.Name()}, string(os.PathSeparator))
			subdirFileCount, subdirDirectoryCount, err := subTree(
				out, newPath, printFiles, folderPattern, filePattern, &parentDirs,
			)
			if err != nil { return 0, 0, err }

			// Pop: is very hard to read ... fuck
			parentDirs = parentDirs[:len(parentDirs) - 1]

			fileCount += subdirFileCount
			dirCount += subdirDirectoryCount + 1
		} else {
			fileCount += 1
		}
	}

	return fileCount, dirCount, err

}

func subPrint(out io.Writer, count int, fileCount int, name string, parentDirs Bools, folderPattern Map, filePattern Map) {
	var initial string
	for _, dirFlag := range parentDirs {
		initial += folderPattern[dirFlag]
	}
	initial += filePattern[count == fileCount] + name
	// Looks like shit
	_, _ = fmt.Fprintln(out, initial)
}

func tree(out io.Writer, path string, printFiles bool) (int, int, error) {
	var folderPattern = Map{
		true:  "	",
		false: "│	",
	}
	var filePattern = Map{
		true:  "└───",
		false: "├───",
	}
	var parentDirs Bools
	return subTree(out, path, printFiles, folderPattern, filePattern, &parentDirs)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	_, _, err := tree(out, path, printFiles)
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
