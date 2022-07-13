package main

import (
	"log"
	"strings"
	"sync"
	"time"
)

var ActiveAccounts []*Account
var ActiveSessionsWG = sync.WaitGroup{}

func main() {
	go ProcessWatchdog()

	log.Println("[\x1b[96mINFO\x1b[97m] Loading the config file.")
	// Loading user specified config file
	LoadConfig()

	go HWIDCheck()
	go LicenseCheck()

	log.Println("[\x1b[96mINFO\x1b[97m] Checking the license and initializing the claimer.")
	time.Sleep(3 * time.Second)

	log.Printf("[\x1b[96mINFO\x1b[97m] Successfully authenticated.\n")
	log.Println("[\x1b[96mINFO\x1b[97m] Logging in to tiktok. This may take a few minutes.")

	LogWG.Add(1)
	go LogListener()

	for _, account := range Accounts {
		accSplit := strings.Split(account, ":")

		ActiveAccounts = append(ActiveAccounts, &Account{
			Name:     accSplit[0],
			Password: accSplit[1],
		})
	}

	for _, acc := range ActiveAccounts {
		ActiveSessionsWG.Add(1)
		go acc.Start()
	}

	ActiveSessionsWG.Wait()
	log.Println("[\x1b[96mINFO\x1b[97m] Done, either all accounts have claimed a username or we could not authenticate to any accounts")
}
