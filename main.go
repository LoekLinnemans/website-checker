package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("couldn't create logfile", err)
		os.Exit(1)
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)

}

func main() {
	var url string
	fmt.Print("Enter the Website: ")
	fmt.Scan(&url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
