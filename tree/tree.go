package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func dirTreeCheckLast(out io.Writer, path string, printFiles bool, prefix string) (err error) {
	// get list of all filenames in current directory
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil
	}

	// remove files from list if nessesary
	for idx := 0; idx < len(list); idx++ {
		if !printFiles && !list[idx].IsDir() {
			list = append(list[:idx], list[idx+1:]...)
			idx--
		}
	}

	// and sort them
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	// for each element in folder
	for idx, file := range list {

		// if its the last file - print "└" instead of "├"
		filePrefix := ""
		if idx == len(list)-1 {
			filePrefix = "└───"
		} else {
			filePrefix = "├───"
		}
		fmt.Fprint(out, prefix+filePrefix+file.Name())

		// if its directory - go inside recursively
		// otherwise print file size
		if file.IsDir() {
			fmt.Fprint(out, "\n")
			var err error
			if file.Name() == list[len(list)-1].Name() {
				err = dirTreeCheckLast(out, path+string(os.PathSeparator)+file.Name(), printFiles, prefix+"\t")
			} else {
				err = dirTreeCheckLast(out, path+string(os.PathSeparator)+file.Name(), printFiles, prefix+"│\t")
			}
			if err != nil {
				panic(err.Error())
			}
		} else if file.Size() > 0 {
			fmt.Fprint(out, " ("+strconv.Itoa(int(file.Size()))+"b)\n")
		} else {
			fmt.Fprint(out, " (empty)\n")
		}

	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	return dirTreeCheckLast(out, path, printFiles, "")
}
