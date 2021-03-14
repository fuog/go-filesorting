package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	go FilePathWalker(Conf.Basics.InputFolder, &Q)

	for {
		time.Sleep(5 * time.Second)
		f, err := Q.get()
		fmt.Println("get", f, err)
		time.Sleep(5 * time.Second)
		if f != nil {
			err2 := Q.remove(*f)
			fmt.Println("remove", err2)
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
