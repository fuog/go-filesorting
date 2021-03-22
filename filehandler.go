package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	Tags        []string
}

func (f *File) GetContentType() error {
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

// checkType returns true if the file matches the requested type.
// Types to check are : pdf
func (f *File) CheckType(expectedType string) (matched bool, reason string) {
	// force type it lowercase
	expectedType = strings.ToLower(expectedType)
	// just make sure ContentType is set
	if f.ContentType == "" {
		if err := f.GetContentType(); err != nil {
			log.Debugln("somthing went wrong calling GetContentType within checkType")
			return false, "ContentType cloud not be optained"
		}

	}
	switch expectedType {
	case "pdf":
		regexmatch1, err := regexp.MatchString(Conf.FileHandling.FileTypePDF.ContentTypeFilter, f.ContentType)
		if err != nil {
			log.Debugln("contenttype does not macht the type pdf")
			return false, "contenttype does not macht"
		}
		regexmatch2, err := regexp.MatchString(Conf.FileHandling.FileTypePDF.FileNameFilter, f.Name)
		if err != nil {
			log.Debugln("filename does not macht the type pdf")
			return false, "filename does not macht"
		}
		// make sure type and filename matches
		return regexmatch1 && regexmatch2, ""
	default:
		log.Warningln("func checkType: no typ maches!, bad scenario")
		return false, ""
	}
}

// readPDF extracts the content of a pdf to f.contentPDF
func (f *File) ReadPdf() error {
	// This is a PDF feature check for type pdf
	if result, _ := f.CheckType("pdf"); !result {
		log.Errorln("This should only be called for pdf files")
		return errors.New("this should only be called for pdf files")
	}
	// make sure the file is not to big
	if f.Size.MBytes() > 10.0 {
		return errors.New("to big to read as PFD")
	}

	openFile, r, err := pdf.Open(f.Path)
	if err != nil {
		log.Errorln(err)
	}
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

func (f *File) DetectTags() error {
	if result, _ := f.CheckType("pdf"); result {

		// iterate over all potenial tags and try to match them
		for _, potentialTag := range Conf.Tagging {
			regexmatchd, err := regexp.MatchString(potentialTag.SearchExpression, f.ContentPDF)
			fmt.Println(potentialTag.Tag)
			if regexmatchd && err == nil {
				log.Debug("the potential tag '" + potentialTag.Tag + "' matches the file '" + f.Name + "'")

				f.Tags = append(f.Tags, potentialTag.Tag)
				f.Tags = append(f.Tags, potentialTag.AdditionalTags...)
				log.Debugf(" adding the following tags : %s, %s", potentialTag.Tag, potentialTag.AdditionalTags)

				return nil
			} else if err != nil {
				// there must be an errors
				return errors.New("error at file: " + f.Name + " and potential Tag : " + potentialTag.Tag + "error is :" + err.Error())
			}

		}

	} else {
		// just in case
		log.Errorln("can not detect tags on a unknown type:", f.Name)
		return errors.New("can not detect tags on a unknown type: " + f.Name)
	}
	return nil
}
