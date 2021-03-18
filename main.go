package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
	log "github.com/sirupsen/logrus"
)

var (
	// Conf Create a new config instance.
	Conf config

	Q FileQueue

	configPath *string
	logPath    *string
)

func init() {

	// Log in default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})

	// setting stdout aslong we do not know the Config
	log.SetOutput(os.Stdout)

	// get configuration for this program
	getConf(&Conf)

	// setting loglevel
	if strings.EqualFold(Conf.Basics.LogLevel, "debug") {
		log.SetLevel(log.DebugLevel)
	} else if strings.EqualFold(Conf.Basics.LogLevel, "info") {
		log.SetLevel(log.InfoLevel)
	} else if strings.EqualFold(Conf.Basics.LogLevel, "warn") {
		log.SetLevel(log.WarnLevel)
	} else if strings.EqualFold(Conf.Basics.LogLevel, "error") {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.DebugLevel)
		log.Warnf("Loglevel \"%s\" is not expected! setting debug loglevel for you", Conf.Basics.LogLevel)
	}

	// based on os args, use file or std out
	// TODO: file usage!
	if Conf.Basics.LogToStdout {
		// Output to stdout instead of the default stderr
		log.SetOutput(os.Stdout)
		log.Debug("Log to stdout is set")
	} else {
		log.Warn("Logfile is set! TODO! using stdout anyway!")
	}

}

func main() {
	fmt.Println("---file read---------")

	// start a goroutine for the filepathwalker
	go FilePathWalker(Conf.Basics.InputFolder, &Q, time.Duration(1E+9*Conf.Basics.ScanInterval))

	for {
		//pdf.DebugOn = true
		time.Sleep(5 * time.Second)
		f, _ := Q.get()
		fmt.Println("file to work with is", f.Name)
		fmt.Println("Q is ", len(Q.files))

		if f != nil {
			//content, err := readPdf(f.Path)
			//fmt.Println("content :", content)

			f.GetFileContentType()
			f.ReadPdf()
			//fmt.Println("content-err :", err)
			//fmt.Println("get", f, err)
			//result, err := regexp.MatchString("(?i)definition", content)
			//fmt.Println("rex-content:", result)
			//fmt.Println("rex-err:", err)
			fmt.Println(f.ContentType, f.Name)
			fmt.Println(f.ContentPDF)
			time.Sleep(2 * time.Second)
			fmt.Println(f)

			err2 := Q.remove(*f)
			if err2 != nil {
				fmt.Println("remove faild:", err2)
			}
			// empty out for new cycle
			f = nil
		}
		fmt.Println("------------")

	}

	//var value int
	//for {
	//	value++
	//	fmt.Println(value)
	//	time.Sleep(time.Second.Round(4))
	//	if value == 10 {
	//		break
	//	}
	//}

}
func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
