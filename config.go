package crongo

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ClientConfig struct {
	Server string
	Token  string
}

func ReadClientConfig(confPath string) ClientConfig {
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("ERROR: Unable to read config: %s", confPath)
		log.Panic()
	}
	var c ClientConfig
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Printf("ERROR: Config is invalid JSON")
		log.Printf("%s\n", err)
		log.Panic()
	}
	return c
}

type ServerConfig struct {
	ValidTokens   []string
	OutputDir     string
	ListenAddress string
}

func ReadServerConfig(confPath string) ServerConfig {
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("ERROR: Unable to read config: %s", confPath)
		log.Panic()
	}
	var sc ServerConfig
	err = json.Unmarshal(b, &sc)

	if err != nil {
		log.Printf("ERROR: Config is invalid JSON")
		log.Printf("%s\n", err)
		log.Panic()
	}
	return sc
}
