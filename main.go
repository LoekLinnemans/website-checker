package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("couldn't create logfile", err)
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)

}

func main() {
	websites, err := Readconfig("config.txt")
	if err != nil {
		log.Fatal("Could not read websites from config.txt:", err)
	}

	for _, url := range websites {
		log.Println(url)
		fmt.Scan(url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
	}
}

func Readconfig(configFile string) ([]string, error) {
	var websites []string
	config, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer config.Close()

	scanner := bufio.NewScanner(config)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			websites = append(websites, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return websites, nil
}
