package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var configPath *string
var logPath *string

func init() {
	getConf()

	// creating arument Parser for cli args
	// argParser := argparse.NewParser("go-filesorter", "A tool for renaming and sorting files")
	// get argument from cli
	// configPath = argParser.String("c", "configfile", &argparse.Options{Required: false, Help: "path to configfile [empty = ./config.yml]"})
	// logPath = argParser.String("l", "logfile", &argparse.Options{Required: false, Help: "path to configfile [empty = stdout ]"})

	//do the actual parsing!
	// err := argParser.Parse(os.Args)
	// if err != nil {
	// 	fmt.Print(argParser.Usage(err))
	// 	os.Exit(1)
	// }

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// Log in default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})
	// based on os args, use file or std out
	// TODO: file usage!
	if Conf.Basics.LogToStdout {
		// Output to stdout instead of the default stderr
		log.SetOutput(os.Stdout)
		log.Debug("Log to stdout is not set, using stdout")
	} else {
		log.Warn("Logfile is set! TODO! using stdout anyway!")
	}

}

func main() {

	fmt.Println("------------")
}
