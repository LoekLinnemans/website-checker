package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("couldn't create logfile", err)
	}
	log.SetOutput(os.Stdout)

	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)

}

func main() {
	websites, err := Readconfig("config.txt")
	if err != nil {
		log.Fatal("Could not read websites from config.txt:", err)
	}

	if err := validateConfig(websites); err != nil {
		log.Fatal("Invalid config file:", err)
	}

	for _, url := range websites {
		if resp, err := http.Get(url); err != nil {
			log.Println("Error:", err)
		} else {
			writeresult(url, resp.StatusCode)
		}
	}

	log.Println("Monitoring complete. Keeping the container alive...")
	time.Sleep(24 * time.Hour)

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

func writeresult(url string, statuscode int) {
	resultfile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Error opening result.txt:", err)
		return
	}
	defer resultfile.Close()

	var result string
	if statuscode == 200 {
		result = fmt.Sprintf("Status code for %s: %d - site reachable\n", url, statuscode)
	} else {
		result = fmt.Sprintf("Status code for %s: %d - site not reachable\n", url, statuscode)
	}

	if _, err := resultfile.WriteString(result); err != nil {
		log.Println("Error writing to result.txt:", err)
	}
}
func validateConfig(websites []string) error {
	for _, url := range websites {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return fmt.Errorf("invalid URL format: %s", url)
		}
		if !strings.Contains(url, ".") {
			return fmt.Errorf("invalid URL format, missing domain: %s", url)
		}
	}
	return nil
}
