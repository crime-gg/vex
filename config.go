package main

import (
	"gopkg.in/ini.v1"
	"log"
	"strings"
	"sync"
)

type ConfigModel struct {
	Username      string `ini:"Username"`
	License       string `ini:"License"`
	AccountsFile  string `ini:"AccountsFile"`
	UsernamesFile string `ini:"UsernamesFile"`
	ProxyFile     string `ini:"ProxyFile"`
	LogDir        string `ini:"LogDirectory"`
	SingleProxy   string `ini:"SingleProxy"`
}

var Config *ConfigModel
var Proxies []string
var ProxyMutex *sync.Mutex
var Accounts []string
var Usernames []string

func LoadConfig() {
	Config = new(ConfigModel)

	var err error

	err = ini.MapTo(Config, "./config.ini")
	if err != nil {
		log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not load the config file. Error: %s\n", err.Error())
	}

	log.Printf("[\x1b[96mINFO\x1b[97m] Successfully loaded the config.")
	if Config.ProxyFile != "null" {
		Proxies = ReadFile(Config.ProxyFile)
		log.Printf("[\x1b[96mINFO\x1b[97m] Successfully loaded %d proxies.\n", len(Proxies))
		for _, proxy := range Proxies {
			if !strings.Contains(proxy, "http") {
				log.Fatalf("[\x1b[91mERROR\x1b[97m] Invalid proxy: %s. Correct syntax: http://IP:PORT(https also works, port is optional, domains are also supported instead of IP)\n",
					proxy)
			}
		}
	}

	Accounts = ReadFile(Config.AccountsFile)
	log.Printf("[\x1b[96mINFO\x1b[97m] Successfully loaded %d accounts.\n", len(Accounts))
	// checking syntax of accounts
	for _, acc := range Accounts {
		accSplit := strings.Split(acc, ":")
		if len(accSplit) != 2 {
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Invalid account: %s. Correct syntax: username:password\n", acc)
		}
	}

	Usernames = ReadFile(Config.UsernamesFile)
	log.Printf("[\x1b[96mINFO\x1b[97m] Successfully loaded %d usernames.\n", len(Usernames))

	ProxyMutex = &sync.Mutex{}
}
