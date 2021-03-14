package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config struct of the entire config
type config struct {
	Basics struct {
		LogToStdout      bool   `mapstructure:",logToStdout"`
		LogLevel         string `mapstructure:",logLevel"`
		LogFile          string `mapstructure:",logFile"`
		InputFolder      string `mapstructure:",inputFolder"`
		ScanInterval     int64  `mapstructure:",scanInterval"`
		FilterdFiles     string `mapstructure:",filterdFiles"`
		Outputfolder     string `mapstructure:",outputfolder"`
		CleanInputFolder struct {
			SortOut          bool   `mapstructure:",SortOut"`
			SortOutFolder    string `mapstructure:",sortOutFolder"`
			SortOutException string `mapstructure:",sortOutException"`
		}
	}
}

func getConf(C *config) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	// reading configfile and create it if needed
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warningln("Configfile not found!")
			// in this case try to create the file
			restoreConf()
			if err := viper.ReadInConfig(); err != nil {
				log.Fatalln("Failed to read config after creating it \n", err)
			}

		} else {
			log.Fatalln("Somthing went wrong reading the config \n", err)
		}
	}

	if err := viper.Unmarshal(C); err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		os.Exit(1)
	}

}

func restoreConf() {
	// template to use if no config is present
	templateString := `
---
basics:
  # log directly to console
  logToStdout: true
  # set loglevel debug|info|warn|fatal
  logLevel: debug
  # only relevant if logToStdout is false
  logFile: ./logfile.log

  # Main Hotfolder that is used
  inputFolder: ./in
  # scanning interval in seconds
  scanInterval: 2
  # the correct filename that is lookd for
  filterdFiles: ".*\\.pdf$"
  # the target folder where to sort files to
  outputfolder: ./out

  cleanInputFolder:
    # should wrong filenames be sorted out (to another Directory)
    sortOut: true
    # where should wrong files be placed?
    sortOutFolder: /tmp/notused
    # should there be exceptions to the sorting out pattern?
    sortOutException: "config\\.yml$"
`
	templateBytes := []byte(templateString)
	err := ioutil.WriteFile("./config.yml", templateBytes, 0660)
	if err != nil {
		log.Fatalln("Error writing configfile to ./ \n", err)
	} else {
		log.Warningln("Created a new configfile from template at ./config.yml")
	}

}
