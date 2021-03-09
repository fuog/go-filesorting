package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config struct of the entire config
type config struct {
	Basics struct {
		LogToStdout      bool   `mapstructure:",logToStdout"`
		LogFile          string `mapstructure:",logFile"`
		InputFolder      string `mapstructure:",inputFolder"`
		FilterdFiles     string `mapstructure:",filterdFiles"`
		Outputfolder     string `mapstructure:",outputfolder"`
		CleanInputFolder struct {
			SortOut          bool   `mapstructure:",SortOut"`
			SortOutFolder    string `mapstructure:",sortOutFolder"`
			SortOutException string `mapstructure:",sortOutException"`
		}
	}
}

var (
	// Conf Create a new config instance.
	Conf *config
)

func getConf() *config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func restoreConf() {
	// template to use if no config is present
	templateString := `
---
basics:
  # log directly to console
  logToStdout: true
  # only relevant if logToStdout is false
  logFile: ./logfile.log

  # Main Hotfolder that is used
  inputFolder: ./in/
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
		log.Errorln("Error writing configfile :", err)
	} else {
		log.Warningln("Created a new configfile from Template at ./config.yml")
	}

}
