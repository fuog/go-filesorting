package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/ledongthuc/pdf"
	log "github.com/sirupsen/logrus"
)

// File for the worker
type File struct {
	Name        string
	Path        string
	Size        datasize.ByteSize
	ModTime     time.Time
	Locked      bool
	ContentType string
	ContentPDF  string
}

func (f *File) GetFileContentType() error {
	// try to open the file
	openFile, err := os.Open(f.Path)
	if err != nil {
		log.Errorln("could not open file: ", err)
		return err
	}
	// remember to close file
	defer openFile.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = openFile.Read(buffer)
	if err != nil {
		log.Errorln("could not read file: ", err)
		return err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	f.ContentType = http.DetectContentType(buffer)

	return nil
}

func (f *File) ReadPdf() error {

	// make sure the file is not to big
	if f.Size.MBytes() > 10.0 {
		return errors.New("to big to read as PFD")
	}

	openFile, r, err := pdf.Open(f.Path)
	// remember to close file
	defer openFile.Close()

	if err != nil {
		return err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return err
	}
	buf.ReadFrom(b)
	f.ContentPDF = buf.String()
	return nil
}
