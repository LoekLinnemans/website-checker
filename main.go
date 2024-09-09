package main

import (
	"fmt"
	"net/http"
)

func main() {
	var url string
	fmt.Print("Enter the Website: ")
	fmt.Scan(&url)
	resp, err := http.Get(url)
	fmt.Print(resp)
	fmt.Print(err)
}
