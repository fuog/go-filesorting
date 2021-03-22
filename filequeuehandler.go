package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
)

// == F I L E Q U E U E =========================

// FileQueue list of files with methods
type FileQueue struct {
	files []File
}

// check if the FileQueue is Empty
func (q *FileQueue) empty() bool {
	if len(q.files) > 0 {
		return false
	}
	return true
}

// add file to filequeue
// TODO: missing errorhandling
func (q *FileQueue) add(f File) {
	// empty check if the queue is empty
	if q.empty() {
		q.files = append(q.files, f)
		log.Infoln("added", f.Name, "to FileQueue")
		return
	}
	// stopped here ! Want to add file is NOT in queue
	for _, entry := range q.files {
		if entry.Path == f.Path {
			// log.Debugln("skipping", f.Name, " file already in queue")
			return
		}
	}
	q.files = append(q.files, f)
	log.Infoln("added", f.Name, "to FileQueue")
	return
}

// list creates a simple list of files in string format for easy debuging
func (q *FileQueue) list() string {

	list := "List of all entries:"
	for _, entry := range q.files {
		list = list + "\n" + entry.Name
	}
	return list
}

// get gets the first not locked file from the FileQueue and locks it before handing over
func (q *FileQueue) get() (*File, error) {
	if q.empty() {
		return nil, errors.New("error: queue empty")
	}
	for i := range q.files {
		if !q.files[i].Locked {
			q.files[i].Locked = true
			var f File
			f = q.files[i]
			return &f, nil
		}
	}
	return nil, errors.New("error: All file in list are locked")
}

func (q *FileQueue) remove(f File) error {

	for i := range q.files {
		if q.files[i].Path == f.Path {
			// TODO: some logging
			// actualy removing in golang .. xD
			q.files = append(q.files[:i], q.files[i+1:]...)
			return nil
		}
	}
	return errors.New("error: No eilequeue entrie matches this file")
}

// == F I L E P A T H W A L K E R ===============

// FilePathWalker is starts a loop that scanns the folder for files and tries to add them to the queue
func FilePathWalker(inputFolder string, q *FileQueue, interval time.Duration) error {
	for {
		// beware the embedded func that will do stuff
		filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
			// this is a loop of all files found

			// this will skip anthing files in subdirectories
			if info.IsDir() && path != inputFolder {
				return filepath.SkipDir
			}
			// ignore folders and unix hidden files
			// TODO: Windows-hidden files work difrently
			if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
				// this is a file we want to use
				// log.Debugln("found a file: " + info.Name())
				s := datasize.ByteSize(info.Size())
				// create new File Object
				f := File{
					Name:    info.Name(),
					Path:    path,
					Size:    s,
					ModTime: info.ModTime(),
					Locked:  false,
				}
				// add to file Queue
				q.add(f)

			}

			return nil
		})

		// the scanning frequency
		time.Sleep(2 * time.Second)
	}
}
