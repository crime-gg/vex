package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

var LogChannel = make(chan string, 100000)

var LogWG = sync.WaitGroup{}

func LogListener() {
	var logMessage string

	var writer *bufio.Writer
	var file *os.File
	var err error

	if Config.LogDir != "" {
		var fileName = Config.LogDir + "/" + time.Now().Format("02_01_2006 15_04_05") + ".log"
		// Creating | Opening log file
		file, err = os.Open(fileName)
		if os.IsNotExist(err) {
			file, err = os.Create(fileName)
			if err != nil {
				log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not create log file! Error: %s\n", err.Error())
			}

		} else if err != nil {
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not open log file! Error: %s\n", err.Error())
		}

		writer = bufio.NewWriter(file)
	}

	if file == nil {
		log.Printf("I AM NIL AS FILE\n")
	}

	if writer == nil {
		log.Printf("I AM NIL AS WRITER\n")
	}

	for {
		logMessage = <-LogChannel

		if logMessage == "EOF" {
			LogWG.Done()
			break
		}

		if Config.LogDir != "" {
			_, err = writer.WriteString(logMessage + "\n")
			if err != nil {
				log.Printf("[WARNING] Could not write to log file. Error: %s\n", err.Error())
			}
			err = writer.Flush()
			if err != nil {
				log.Printf("[WARNING] Could not write to log file. Error: %s\n", err.Error())
			}
		}
	}
}
