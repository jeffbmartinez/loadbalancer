package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jeffbmartinez/loadbalancer/host"
)

type Config struct {
	ListenLocalOnly bool        `json:"listenLocalOnly"`
	ListenPort      int         `json:"listenPort"`
	Hosts           []host.Host `json:"hosts"`
}

/*
NewConfig takes a filename and returns a configuration object.
*/
func NewConfig(filename string) (conf Config, err error) {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	json.Unmarshal(fileContents, &conf)

	err = conf.verifyConfig()

	return conf, err
}

func (c Config) ListenAddress() string {
	listenHostname := ""
	if c.ListenLocalOnly {
		listenHostname = "localhost"
	}

	return fmt.Sprintf("%v:%v", listenHostname, c.ListenPort)
}

/*
Display prints the contents of the configuration in a friendly way.
*/
func (c Config) Display() {
	listenAddresses := "Accepting all connections"
	if c.ListenLocalOnly {
		listenAddresses = "Accepting connections from localhost only (no remote connections allowed)"
	}

	fmt.Println(listenAddresses)
	fmt.Printf("Listening on port %v\n\n", c.ListenPort)

	fmt.Printf("Balancing across %v hosts:\n", len(c.Hosts))
	for _, host := range c.Hosts {
		fmt.Printf("\tHostname: '%v', weight: %v\n", host.Hostname, host.Weight)
	}
}

/*
verifyConfig checks for possible issues with loaded configuration settings.
*/
func (c Config) verifyConfig() error {
	const minListenPort int = 1
	const maxListenPort int = 65535

	if c.ListenPort < minListenPort || c.ListenPort > maxListenPort {
		return errors.New("listenPort is not within valid range")
	}

	if len(c.Hosts) == 0 {
		return errors.New("List of hosts is empty")
	}

	for _, host := range c.Hosts {
		if host.Weight <= 0 {
			return errors.New("Found invalid weight for at least one host")
		}
	}

	return nil
}
