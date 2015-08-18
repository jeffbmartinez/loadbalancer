package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/loadbalancer/config"
	"github.com/jeffbmartinez/loadbalancer/handler"
)

const exitFailure = 1
const exitUsageError = 2

func main() {
	configFilename := getConfigFilename()

	conf, err := config.NewConfig(configFilename)
	if err != nil {
		log.Errorf("Trouble reading config file '%v' (%v)\n", configFilename, err)
		os.Exit(exitFailure)
	}

	log.Printf("Using config file '%v'\n\n", configFilename)
	conf.Display()

	http.Handle("/", handler.NewBalancer(conf.Hosts))

	http.ListenAndServe(conf.ListenAddress(), nil)
}

func getConfigFilename() (configFilename string) {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Errorf("Usage: loadbalancer path/to/config.json")
		os.Exit(exitUsageError)
	}

	return flag.Arg(0)
}
