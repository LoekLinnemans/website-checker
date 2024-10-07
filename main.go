package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("couldn't create logfile: %v", err)
	}
	defer logFile.Close()

	resultfile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("couldn't create resultfile: %v", err)
	}
	defer resultfile.Close()

	configfile, err := os.OpenFile("config.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("couldn't create configfile: %v", err)
	}
	defer configfile.Close()

	log.SetOutput(os.Stdout)

	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)

}

func main() {
	for {
		websites, err := Readconfig("config.txt")
		if err != nil {
			log.Fatalf("Could not read websites from config.txt: %v", err)
		}

		if err := validateConfig(websites); err != nil {
			log.Fatalf("Invalid config file: %v", err)
		}

		for _, url := range websites {
			if resp, err := http.Get(url); err != nil {
				log.Printf("Error reaching %s: %v", url, err)
				writeresult(url, 0)
			} else {
				writeresult(url, resp.StatusCode)
			}
		}

		log.Printf("Monitoring complete. Restarting in 1 minutes...")
		time.Sleep(1 * time.Minute)
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

func writeresult(url string, statuscode int) {
	resultfile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening result.txt: %v", err)
	}
	var result string
	if statuscode == 200 {
		result = fmt.Sprintf("Status code for %s: %d - site reachable\n", url, statuscode)
	} else if statuscode == 0 {
		result = fmt.Sprintf("Error reaching %s - host not reachable\n", url)
	} else {
		result = fmt.Sprintf("Status code for %s: %d - site not reachable\n", url, statuscode)
	}

	if _, err := resultfile.WriteString(result); err != nil {
		log.Printf("Error writing to result.txt: %v", err)
	}
}

func validateConfig(websites []string) error {
	var urlRegex = regexp.MustCompile(`^(http|https)://www\.[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)

	for _, url := range websites {
		if !urlRegex.MatchString(url) {
			return fmt.Errorf("invalid URL format: %s", url)
		}
	}
	return nil
}
