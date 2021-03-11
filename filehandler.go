package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// File for the worker
type File struct {
	file   os.FileInfo
	Locked bool
}

// FileQueue list of files with methods
type FileQueue []File

// check if the FileQueue is Empty
func (q *FileQueue) empty() bool {
	if len(*q) > 0 {
		return false
	}
	return true
}

// add file to filequeue
// TODO: missing errorhandling
func (q *FileQueue) add(f File) {
	// empty check if the queue is empty
	if q.empty() {
		*q = append(*q, f)
		log.Infoln("added", f.file.Name(), "to FileQueue")
		return
	}
	// stopped here ! Want to add file is NOT in queue
	for _, entry := range *q {
		if entry == f {
			log.Infoln("skipping", f)
			return
		}
	}
	*q = append(*q, f)
	log.Infoln("added", f)
	return
}

func (q *FileQueue) list() string {

	list := "List of all entries: \n"

	for _, entry := range *q {
		list = list + entry.file.Name() + "\n"
	}
	return list
}

// FilePathWalkDir TODO: writing this Comment
func FilePathWalkDir(inputFolder string, q *FileQueue) error {
	err := filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
		// this is a loop of all files found

		// ignore unix hidden files
		// TODO: Windows works difrent
		if info.IsDir() && path != inputFolder {
			return filepath.SkipDir
		}

		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			// this is a file we want to use
			fmt.Println("found a file: " + info.Name())
			// create new File Object
			f := File{info, false}
			// add to file Queue
			q.add(f)

		}

		return nil
	})
	return err
}
