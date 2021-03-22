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

	// everything very basic
	Basics struct {
		// enable stdout
		LogToStdout bool `mapstructure:",logToStdout"`
		// Mode [debug/info/warn/error]
		LogLevel string `mapstructure:",logLevel"`
		// path
		LogFile string `mapstructure:",logFile"`
	}
	// everything about handling files
	FileHandling struct {
		// path
		InputFolder string `mapstructure:",inputFolder"`
		// intervall in seconds
		ScanInterval int64 `mapstructure:",scanInterval"`
		// path
		SortOutFolder string `mapstructure:",sortOutFolder"`
		// regexp
		IgnoredFileNames string `mapstructure:",ignoredFileNames"`

		// everything about identifying PDF files
		FileTypePDF struct {
			ContentTypeFilter string `mapstructure:",contentTypeFilter"`
			FileNameFilter    string `mapstructure:",fileNameFilter"`
		}
	}
	// List of all tags-filters
	Tagging []taggingConf
}

type taggingConf struct {
	Tag              string   `mapstructure:",tag"`
	SearchExpression string   `mapstructure:",searchExpression"`
	AdditionalTags   []string `mapstructure:",additionalTags"`
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
	templateString := `---
basics:
  # log directly to console
  logToStdout: true
  # set loglevel debug|info|warn|fatal
  logLevel: debug
  # only relevant if logToStdout is false
  logFile: ./logfile.log

fileHandling:
  # Main Hotfolder that is used
  inputFolder: ./in
  # scanning interval in seconds
  scanInterval: 10
  # where should wrong files be placed?
  sortOutFolder: /tmp/notused
  # should there be exceptions to the sorting out pattern?
  ignoredFileNames: "config\\.yml$"

  # passible fileTypes
  fileTypePDF:
	contentTypeFilter: application/pdf
	fileNameFilter: ".*\\.pdf$"

# List can be extendet at will but with th (correct types!)
tagging:
  - tag: swisscom
	searchExpression: "swisscom"
	additionaltags:
	  - isp
  - tag: cablecom
	searchExpression: "cablecom"
	additionaltags:
	  - isp
`
	templateBytes := []byte(templateString)
	err := ioutil.WriteFile("./config.yml", templateBytes, 0660)
	if err != nil {
		log.Fatalln("Error writing configfile to ./ \n", err)
	} else {
		log.Warningln("Created a new configfile from template at ./config.yml")
	}

}
