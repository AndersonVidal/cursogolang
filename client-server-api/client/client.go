package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const quotationEndpoint = "http://localhost:8080/cotacao"
const quotationsFile = "cotacao.txt"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", quotationEndpoint, nil)
	if err != nil {
		log.Fatal("Error in make get quotation request: " + err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error in get quotation: " + err.Error())
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code in get quotation response: %d", res.StatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	registerFileQuotation(string(body))
}

func registerFileQuotation(quotation string) {
	var file *os.File
	_, err := os.Stat(quotationsFile)
	if os.IsNotExist(err) {
		file, err = os.Create(quotationsFile)
		if err != nil {
			log.Fatal("Error to create quotations file: ", err)
		}
	} else {
		file, err = os.OpenFile(quotationsFile, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatal("Error to open quotations file: ", err)
		}
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + strings.Trim(strings.TrimSpace(quotation), `"`) + "\n")
	if err != nil {
		log.Fatal("Error to add quotation in quotations file: ", err)
	}
}
